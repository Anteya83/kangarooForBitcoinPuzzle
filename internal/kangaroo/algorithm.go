package kangaroo

import (
	"crypto/sha256"
	"fmt"
	"kangarooForBitcoinPuzzle/internal/config"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"kangarooForBitcoinPuzzle/internal/shardedmap"
	"kangarooForBitcoinPuzzle/internal/state"
	"kangarooForBitcoinPuzzle/internal/types"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type KangarooSolver struct {
	cfg               config.Config
	a, b              *big.Int
	pubKey            *types.Point
	tameMap           *shardedmap.ShardedMap
	tamePos           *types.Point
	tameDist          *big.Int
	found             int32
	result            *big.Int
	mu                sync.Mutex
	stopChan          chan struct{}
	wg                sync.WaitGroup
	stats             *Stats
	jumpTable         []*big.Int
	stateFile         string
	saveCounter       int64
	maxJump           int64
	totalTameJumps    int64
	closeOnce         sync.Once
	posMu             sync.RWMutex
	distinguishedMask *big.Int
	saveMu            sync.Mutex
	jumpPoints        []*types.Point
}
type Stats struct {
	sync.Mutex
	TotalTameSteps  int64
	TotalWildSteps  int64
	StartTime       time.Time
	CollisionTime   time.Time
	TameMapSize     int
	CollisionsFound int32
}

func NewKangarooSolver(cfg config.Config, publicKey *types.Point, a, b *big.Int) *KangarooSolver {
	var effectiveMaxJump int64
	var bits int
	if cfg.MaxJump > 0 {
		effectiveMaxJump = cfg.MaxJump
		width := new(big.Int).Sub(b, a)
		bits = width.BitLen()
	} else {
		effectiveMaxJump, bits = optimalMaxJump(a, b)
	}
	fmt.Printf("Using MaxJump = %d (width bits=%d)\n", effectiveMaxJump, bits)
	jumpTable := make([]*big.Int, effectiveMaxJump)
	for i := int64(1); i <= effectiveMaxJump; i++ {
		jumpTable[i-1] = new(big.Int).Exp(big.NewInt(2), big.NewInt(i-1), nil)
	}
	jumpPoints := make([]*types.Point, effectiveMaxJump)
	for i := int64(1); i <= effectiveMaxJump; i++ {
		jumpPoints[i-1] = crypto.ScalarMult(
			new(big.Int).Exp(big.NewInt(2), big.NewInt(i-1), nil),
			&types.Point{X: crypto.Gx, Y: crypto.Gy},
		)
	}
	mask := new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), uint(cfg.DistinguishedBits)),
		big.NewInt(1),
	)
	stateFile := generateStateFilename(publicKey, a, b)
	return &KangarooSolver{
		cfg:               cfg,
		a:                 new(big.Int).Set(a),
		b:                 new(big.Int).Set(b),
		pubKey:            publicKey,
		tameMap:           shardedmap.NewShardedMap(cfg.NumShards),
		stopChan:          make(chan struct{}),
		stats:             &Stats{StartTime: time.Now()},
		jumpTable:         jumpTable,
		jumpPoints:        jumpPoints,
		maxJump:           effectiveMaxJump,
		distinguishedMask: mask,
		stateFile:         stateFile,
	}
}
func generateStateFilename(pubKey *types.Point, a, b *big.Int) string {
	data := pubKey.X.Text(16) + "|" + pubKey.Y.Text(16) + "|" + a.Text(16) + "|" + b.Text(16)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("state_%x.bin", hash[:8])
}
func (k *KangarooSolver) isDistinguished(p *types.Point) bool {
	var low big.Int
	low.And(p.X, k.distinguishedMask)
	return low.Sign() == 0
}
func (k *KangarooSolver) jumpDistance(p *types.Point) int64 {
	xb := p.X.Bytes()
	yb := p.Y.Bytes()
	h := uint64(14695981039346656037)
	for _, b := range xb {
		h ^= uint64(b)
		h *= 1099511628211
	}
	for _, b := range yb {
		h ^= uint64(b)
		h *= 1099511628211
	}
	idx := int64(h % uint64(k.maxJump))
	return idx
}

func (k *KangarooSolver) saveState(foundKey *big.Int) error {
	if atomic.LoadInt32(&k.found) == 1 {
		k.mu.Lock()
		if k.result != nil {
			foundKey = k.result
		}
		k.mu.Unlock()
	}
	mapData := make(map[string]string)
	k.tameMap.ForEach(func(key string, val *big.Int) {
		mapData[key] = val.Text(16)
	})

	k.posMu.RLock()
	var tameX, tameY, tameDist string
	if k.tamePos != nil && k.tameDist != nil {
		tameX = k.tamePos.X.Text(16)
		tameY = k.tamePos.Y.Text(16)
		tameDist = k.tameDist.Text(16)
	} else {
		tameX, tameY, tameDist = "", "", ""
	}
	k.posMu.RUnlock()

	st := &state.State{
		Version:        1,
		PublicKeyHex:   k.cfg.PublicKeyHex,
		StartRangeHex:  k.cfg.StartRangeHex,
		EndRangeHex:    k.cfg.EndRangeHex,
		TamePosX:       tameX,
		TamePosY:       tameY,
		TameDist:       tameDist,
		TotalTameSteps: atomic.LoadInt64(&k.stats.TotalTameSteps),
		TotalWildSteps: atomic.LoadInt64(&k.stats.TotalWildSteps),
		Map:            mapData,
	}
	if foundKey != nil {
		st.FoundKey = foundKey.Text(16)
	}
	k.saveMu.Lock()
	defer k.saveMu.Unlock()
	return st.Save(k.stateFile)
}
func (k *KangarooSolver) loadState() (bool, *big.Int, error) {
	st, err := state.Load(k.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("failed to load state: %w", err)
	}
	if !st.IsValid(k.cfg.PublicKeyHex, k.cfg.StartRangeHex, k.cfg.EndRangeHex) {

		return false, nil, nil
	}
	if st.FoundKey != "" {
		fk := new(big.Int)
		if _, ok := fk.SetString(st.FoundKey, 16); !ok {
			return false, nil, fmt.Errorf("invalid FoundKey: %s", st.FoundKey)
		}
		return true, fk, nil
	}
	if st.TamePosX == "" || st.TamePosY == "" || st.TameDist == "" {
		return false, nil, fmt.Errorf("state has empty coordinates, ignoring")
	}
	tameX, ok := new(big.Int).SetString(st.TamePosX, 16)
	if !ok {
		return false, nil, fmt.Errorf("invalid TamePosX: %s", st.TamePosX)
	}
	tameY, ok := new(big.Int).SetString(st.TamePosY, 16)
	if !ok {
		return false, nil, fmt.Errorf("invalid TamePosY: %s", st.TamePosY)
	}
	tameDist, ok := new(big.Int).SetString(st.TameDist, 16)
	if !ok {
		return false, nil, fmt.Errorf("invalid TameDist: %s", st.TameDist)
	}

	k.posMu.Lock()
	k.tamePos = types.NewPoint(tameX, tameY)
	k.tameDist = tameDist
	k.posMu.Unlock()
	for key, distStr := range st.Map {
		dist, _ := new(big.Int).SetString(distStr, 16)
		if dist != nil {
			k.tameMap.Set(key, dist)
		}
	}
	atomic.StoreInt64(&k.stats.TotalTameSteps, st.TotalTameSteps)
	atomic.StoreInt64(&k.stats.TotalWildSteps, st.TotalWildSteps)
	return true, nil, nil
}
func (k *KangarooSolver) tameKangaroo() {
	defer k.wg.Done()
	k.posMu.Lock()
	if k.tamePos == nil {
		k.tamePos = crypto.ScalarMult(k.b, &types.Point{X: crypto.Gx, Y: crypto.Gy})
		k.tameDist = big.NewInt(0)
	}
	pos := &types.Point{
		X: new(big.Int).Set(k.tamePos.X),
		Y: new(big.Int).Set(k.tamePos.Y),
	}
	dist := new(big.Int).Set(k.tameDist)
	k.posMu.Unlock()

	saveCounter := int64(0)
	totalJumps := int64(0)
	for {
		select {
		case <-k.stopChan:
			if atomic.LoadInt32(&k.found) == 1 {
				k.mu.Lock()
				res := k.result
				k.mu.Unlock()
				k.saveState(res)
			} else {
				k.saveState(nil)
			}
			return
		default:
			if totalJumps > k.cfg.TameStepsLimit {
				k.saveState(nil)
				k.closeOnce.Do(func() { close(k.stopChan) })
				return
			}
			if k.isDistinguished(pos) {
				key := fmt.Sprintf("%x,%x", pos.X, pos.Y)
				k.tameMap.Set(key, new(big.Int).Set(dist))
				atomic.AddInt64(&k.stats.TotalTameSteps, 1)
				saveCounter++
				if saveCounter%int64(k.cfg.SaveInterval) == 0 {
					if atomic.LoadInt32(&k.found) == 1 {
						k.mu.Lock()
						res := k.result
						k.mu.Unlock()
						k.saveState(res)
					} else {
						k.saveState(nil)
					}
				}
			}

			idx := k.jumpDistance(pos)
			jump := k.jumpTable[idx]
			jumpPoint := k.jumpPoints[idx]
			pos = crypto.AddPoints(pos, jumpPoint)
			dist.Add(dist, jump)
			totalJumps++

			k.posMu.Lock()
			k.tamePos = pos
			k.tameDist = new(big.Int).Set(dist)
			k.posMu.Unlock()
		}
	}
}
func (k *KangarooSolver) wildKangaroo(workerID int) {
	defer k.wg.Done()
	width := new(big.Int).Sub(k.b, k.a)
	step := new(big.Int).Div(width, big.NewInt(int64(k.cfg.NumWorkers*2)))
	if step.Sign() == 0 {
		step = big.NewInt(1)
	}
	offset := new(big.Int).Mul(big.NewInt(int64(workerID+1)), step)
	offset.Mod(offset, width)
	wildPos := crypto.AddPoints(k.pubKey, crypto.ScalarMult(offset, &types.Point{X: crypto.Gx, Y: crypto.Gy}))
	wildDist := new(big.Int).Set(offset)
	for {
		select {
		case <-k.stopChan:
			return
		default:
			if atomic.LoadInt32(&k.found) == 1 {
				return
			}
			if k.isDistinguished(wildPos) {
				key := fmt.Sprintf("%x,%x", wildPos.X, wildPos.Y)
				if tameDist, exists := k.tameMap.Get(key); exists {
					result := new(big.Int).Add(k.b, tameDist)
					result.Sub(result, wildDist)
					result.Mod(result, crypto.N)

					if result.Cmp(k.a) >= 0 && result.Cmp(k.b) <= 0 {
						checkPub := crypto.PubKeyFromPrivate(result)
						if checkPub.X.Cmp(k.pubKey.X) == 0 && checkPub.Y.Cmp(k.pubKey.Y) == 0 {
							if atomic.CompareAndSwapInt32(&k.found, 0, 1) {
								k.mu.Lock()
								k.result = result
								k.mu.Unlock()
								k.stats.Lock()
								k.stats.CollisionTime = time.Now()
								k.stats.Unlock()
								atomic.AddInt32(&k.stats.CollisionsFound, 1)
								k.saveState(result)
								k.closeOnce.Do(func() { close(k.stopChan) })
								return
							}
						}
					}
				}
			}
			idx := k.jumpDistance(wildPos)
			jump := k.jumpTable[idx]
			jumpPoint := k.jumpPoints[idx]
			wildPos = crypto.AddPoints(wildPos, jumpPoint)
			wildDist.Add(wildDist, jump)
			atomic.AddInt64(&k.stats.TotalWildSteps, 1)
		}
	}
}
func (k *KangarooSolver) Solve() *big.Int {
	fmt.Printf("Starting Kangaroo search (Distinguished bits: %d)\n", k.cfg.DistinguishedBits)
	fmt.Printf("Workers: %d\n", k.cfg.NumWorkers)
	fmt.Printf("Range: %x - %x (HEX)\n", k.a, k.b)
	fmt.Printf("State file: %s\n", k.stateFile)

	loaded, foundKey, err := k.loadState()
	if err != nil {
		fmt.Printf("State file is corrupted: %v, removing and starting fresh\n", err)
		os.Remove(k.stateFile)
		loaded = false
	}
	if loaded && foundKey != nil {
		fmt.Printf("Key already found!\n")
		return foundKey
	}
	if loaded {
		fmt.Println("Continuing from saved state...")
	} else {
		fmt.Println("Starting fresh search...")
	}
	done := make(chan struct{})
	defer close(done)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\nInterrupt received, saving state...")
			k.saveState(nil)
			k.closeOnce.Do(func() { close(k.stopChan) })
		case <-done:
		}
	}()

	minDP := 100
	k.wg.Add(1)
	go k.tameKangaroo()
	fmt.Println("Populating distinguished points...")

	for k.tameMap.Len() < minDP {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-k.stopChan:
			k.wg.Wait()
			k.mu.Lock()
			res := k.result
			k.mu.Unlock()
			return res
		default:
		}
	}

	for i := 0; i < k.cfg.NumWorkers; i++ {
		k.wg.Add(1)
		go k.wildKangaroo(i)
	}
	go k.monitorProgress()
	k.wg.Wait()
	k.mu.Lock()
	res := k.result
	k.mu.Unlock()
	return res
}
func (k *KangarooSolver) monitorProgress() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-k.stopChan:
			return
		case <-ticker.C:
			tame := atomic.LoadInt64(&k.stats.TotalTameSteps)
			wild := atomic.LoadInt64(&k.stats.TotalWildSteps)
			size := k.tameMap.Len()
			fmt.Printf("Pgs: Tame: %d, Wild: %d, Map: %d points\n", tame, wild, size)
		}
	}
}

func (k *KangarooSolver) PrintStats() {
	k.stats.Lock()
	defer k.stats.Unlock()
	fmt.Printf("Runtime: %v\n", time.Since(k.stats.StartTime))
	if !k.stats.CollisionTime.IsZero() {
		fmt.Printf("Collision time: %v\n", k.stats.CollisionTime.Sub(k.stats.StartTime))
	}
	fmt.Printf("Tame steps: %d\n", atomic.LoadInt64(&k.stats.TotalTameSteps))
	fmt.Printf("Wild steps: %d\n", atomic.LoadInt64(&k.stats.TotalWildSteps))
	fmt.Printf("Collisions found: %d\n", atomic.LoadInt32(&k.stats.CollisionsFound))
}

func optimalMaxJump(a, b *big.Int) (int64, int) {
	width := new(big.Int).Sub(b, a)
	bits := width.BitLen()
	if bits < 16 {
		return 4, bits
	}
	jump := int64(bits/2 + 1)
	if jump < 4 {
		jump = 4
	}

	return jump, bits
}
