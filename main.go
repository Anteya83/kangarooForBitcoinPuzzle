package main

import (
	"fmt"
	"kangarooForBitcoinPuzzle/internal/kangaroo"
	"kangarooForBitcoinPuzzle/internal/types"
	"math/big"
)

func main() {
	fmt.Println("--Kangaroo start--")

	pubKey := types.NewPoint(big.NewInt(1), big.NewInt(2))
	a := big.NewInt(0)
	b := big.NewInt(100)

	solver := kangaroo.NewKangarooSolver(pubKey, a, b)
	result := solver.Solve()
	fmt.Printf("Result: %v\n", result)
	solver.PrintStats()
}
