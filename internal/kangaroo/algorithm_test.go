package kangaroo

import (
	"kangarooForBitcoinPuzzle/internal/config"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"math/big"
	"os"
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
	defer os.Remove(solver.stateFile)
	found := solver.Solve()
	if found == nil {
		t.Fatal("Solve returned nil")
	}
	if found.Cmp(priv) != 0 {
		t.Errorf("Found %v, expected %v", found, priv)
	}
}
func TestResume(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.NumWorkers = 2
	cfg.TameStepsLimit = 500
	cfg.DistinguishedBits = 4
	cfg.SaveInterval = 10

	priv := big.NewInt(12345)
	pubKey := crypto.PubKeyFromPrivate(priv)
	a := big.NewInt(10000)
	b := big.NewInt(20000)

	solver1 := NewKangarooSolver(cfg, pubKey, a, b)
	defer os.Remove(solver1.stateFile)
	found1 := solver1.Solve()
	if found1 != nil {
		t.Logf("First solve unexpectedly found key: %v (state will be saved anyway)", found1)
	}

	cfg2 := cfg
	cfg2.TameStepsLimit = 1000000
	solver2 := NewKangarooSolver(cfg2, pubKey, a, b)
	defer os.Remove(solver2.stateFile)
	found2 := solver2.Solve()

	if found2 == nil {
		t.Fatal("Resume solve returned nil")
	}
	if found2.Cmp(priv) != 0 {
		t.Errorf("Resume found %v, expected %v", found2, priv)
	}

	checkPub := crypto.PubKeyFromPrivate(found2)
	if checkPub.X.Cmp(pubKey.X) != 0 || checkPub.Y.Cmp(pubKey.Y) != 0 {
		t.Errorf("Resume found key does not match public key")
	}

	os.Remove(solver2.stateFile)
}
