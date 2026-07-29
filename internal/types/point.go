package types

import "math/big"

type Point struct {
	X, Y *big.Int
}

func NewPoint(x, y *big.Int) *Point {
	if x == nil {
		x = new(big.Int)
	}
	if y == nil {
		y = new(big.Int)
	}
	return &Point{X: x, Y: y}
}

func ZeroPoint() *Point {
	return NewPoint(new(big.Int).SetInt64(0), new(big.Int).SetInt64(0))
}
