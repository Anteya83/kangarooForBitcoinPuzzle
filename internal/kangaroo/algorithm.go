package kangaroo

import (
	"kangarooForBitcoinPuzzle/internal/types"
	"math/big"
)

type KangarooSolver struct {
}

func NewKangarooSolver(publicKey *types.Point, a, b *big.Int) *KangarooSolver {
	return &KangarooSolver{}
}

func (ks *KangarooSolver) Solve() *big.Int {
	return nil
}
func (ks *KangarooSolver) PrintStats() {
	
}
