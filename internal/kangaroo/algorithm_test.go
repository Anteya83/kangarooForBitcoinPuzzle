package kangaroo

import (
	"kangarooForBitcoinPuzzle/internal/config"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"math/big"
	"testing"
)

func TestKangarooSolveSmallRange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.NumWorkers = 4
	cfg.TameStepsLimit = 1000000
	cfg.DistinguishedBits = 4
	cfg.SaveInterval = 100
	cfg.MaxJump = 0

	priv := big.NewInt(12345)
	pubKey := crypto.PubKeyFromPrivate(priv)
	a := big.NewInt(10000)
	b := big.NewInt(20000)

	solver := NewKangarooSolver(cfg, pubKey, a, b)
	found := solver.Solve()
	if found == nil {
		t.Fatal("Solve returned nil")
	}
	if found.Cmp(priv) != 0 {
		t.Errorf("Found %v, expected %v", found, priv)
	}
}
