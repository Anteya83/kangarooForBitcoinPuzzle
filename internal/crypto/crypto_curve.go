package crypto

import (
	"encoding/hex"
	"fmt"
	"kangarooForBitcoinPuzzle/internal/types"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

var (
	curve     = btcec.S256()
	N         = curve.Params().N
	Gx, Gy    = curve.Params().Gx, curve.Params().Gy
	zeroPoint = types.ZeroPoint()
)

func ScalarMult(k *big.Int, p *types.Point) *types.Point {
	if k.Cmp(big.NewInt(0)) == 0 {
		return types.ZeroPoint()
	}
	x, y := curve.ScalarMult(p.X, p.Y, k.Bytes())
	return types.NewPoint(x, y)
}

func AddPoints(a, b *types.Point) *types.Point {
	if a.X.Cmp(zeroPoint.X) == 0 && a.Y.Cmp(zeroPoint.Y) == 0 {
		return types.NewPoint(b.X, b.Y)
	}
	if b.X.Cmp(zeroPoint.X) == 0 && b.Y.Cmp(zeroPoint.Y) == 0 {
		return types.NewPoint(a.X, a.Y)
	}
	x, y := curve.Add(a.X, a.Y, b.X, b.Y)
	return types.NewPoint(x, y)
}

func PubKeyFromPrivate(priv *big.Int) *types.Point {
	return ScalarMult(priv, &types.Point{X: Gx, Y: Gy})
}

func UncompressPublicKey(compressedHex string) (x, y *big.Int, err error) {
	compressedBytes, err := hex.DecodeString(compressedHex)
	if err != nil {
		return nil, nil, err
	}
	if len(compressedBytes) != 33 {
		return nil, nil, fmt.Errorf("inval key->len: %d", len(compressedBytes))
	}
	pubKey, err := btcec.ParsePubKey(compressedBytes)
	if err != nil {
		return nil, nil, err
	}
	return pubKey.X(), pubKey.Y(), nil
}

func ParsePubKey(hexKey string) (*types.Point, error) {
	if len(hexKey) > 2 && hexKey[:2] == "0x" {
		hexKey = hexKey[2:]
	}
	bytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	switch len(bytes) {
	case 33:
		x, y, err := UncompressPublicKey(hexKey)
		if err != nil {
			return nil, err
		}
		return types.NewPoint(x, y), nil
	case 65:
		if bytes[0] != 0x04 {
			return nil, fmt.Errorf("inval uncompr key->pref")
		}
		x := new(big.Int).SetBytes(bytes[1:33])
		y := new(big.Int).SetBytes(bytes[33:65])
		return types.NewPoint(x, y), nil
	default:
		return nil, fmt.Errorf("inval key->len: %d bytes", len(bytes))
	}
}
