package main

import (
	"fmt"
	"kangarooForBitcoinPuzzle/internal/cli"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"kangarooForBitcoinPuzzle/internal/encoding"
	"kangarooForBitcoinPuzzle/internal/kangaroo"
	"math/big"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := cli.ParseFlags()

	fmt.Printf("Public key: %s\n", cfg.PublicKeyHex)
	fmt.Printf("Range: %s - %s\n", cfg.StartRangeHex, cfg.EndRangeHex)
	fmt.Printf("Distinguished bits: %d\n", cfg.DistinguishedBits)
	fmt.Printf("Workers: %d\n", cfg.NumWorkers)
	fmt.Printf("MaxJump: %d (0 = auto)\n", cfg.MaxJump)
	fmt.Printf("Tame limit: %d steps\n", cfg.TameStepsLimit)
	pubKey, err := crypto.ParsePubKey(cfg.PublicKeyHex)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	a, b := new(big.Int), new(big.Int)
	a.SetString(cfg.StartRangeHex, 16)
	b.SetString(cfg.EndRangeHex, 16)

	solver := kangaroo.NewKangarooSolver(cfg, pubKey, a, b)
	start := time.Now()
	found := solver.Solve()
	elapsed := time.Since(start)

	solver.PrintStats()
	if found != nil {
		privHex := fmt.Sprintf("%064x", found)
		wif, _ := encoding.PrivateKeyToWIF(privHex)
		pubBytes, _ := encoding.GetPublicKey(privHex)
		compressed := encoding.CompressPublicKey(pubBytes)
		addr, _ := encoding.PublicKeyToAddress(compressed)
		p2sh, _ := encoding.PublicKeyToP2SHAddress(compressed)

		fmt.Printf("PrivKey HEX: %s\n", privHex)
		fmt.Printf("WIF(uncompressed): %s\n", wif)
		fmt.Printf("Addr(P2PKH): %s\n", addr)
		fmt.Printf("Addr(P2SH): %s\n", p2sh)
	} else {
		fmt.Println("Key not found")
	}
	fmt.Printf("Time: %v\n", elapsed)
}
