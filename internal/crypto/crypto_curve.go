package crypto

import (
	"kangarooForBitcoinPuzzle/internal/types"
	"math/big"
)

func ScalarMult(k *big.Int, p *types.Point) *types.Point {
	return types.ZeroPoint()
}

func AddPoints(a, b *types.Point) *types.Point {
	return types.ZeroPoint()
}

func PubKeyFromPrivate(priv *big.Int) *types.Point {
	return types.ZeroPoint()
}
func ParsePubKey(hexKey string) (*types.Point, error) {
	return types.ZeroPoint(), nil
}
