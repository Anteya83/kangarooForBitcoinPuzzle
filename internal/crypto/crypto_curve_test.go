package crypto

import (
	"math/big"
	"testing"

	"kangarooForBitcoinPuzzle/internal/types"
)

func TestAddPoints(t *testing.T) {
	G := &types.Point{X: Gx, Y: Gy}

	// -G = (Gx, -Gy mod P)
	p := curve.Params().P
	negY := new(big.Int).Sub(p, Gy)
	negG := &types.Point{X: new(big.Int).Set(Gx), Y: negY}

	result := AddPoints(G, negG)
	expected := types.ZeroPoint()
	if result.X.Cmp(expected.X) != 0 || result.Y.Cmp(expected.Y) != 0 {
		t.Errorf("G + (-G) != ZeroPoint, got (%x, %x)", result.X, result.Y)
	}

	// G + ZeroPoint = G
	result2 := AddPoints(G, types.ZeroPoint())
	if result2.X.Cmp(G.X) != 0 || result2.Y.Cmp(G.Y) != 0 {
		t.Errorf("G + ZeroPoint != G, got (%x, %x)", result2.X, result2.Y)
	}

	// ZeroPoint + G = G
	result3 := AddPoints(types.ZeroPoint(), G)
	if result3.X.Cmp(G.X) != 0 || result3.Y.Cmp(G.Y) != 0 {
		t.Errorf("ZeroPoint + G != G, got (%x, %x)", result3.X, result3.Y)
	}
}

func TestScalarMult(t *testing.T) {
	G := &types.Point{X: Gx, Y: Gy}

	// 1 * G = G
	result := ScalarMult(big.NewInt(1), G)
	if result.X.Cmp(G.X) != 0 || result.Y.Cmp(G.Y) != 0 {
		t.Errorf("1*G != G, got (%x, %x)", result.X, result.Y)
	}

	// 0 * G = ZeroPoint
	result2 := ScalarMult(big.NewInt(0), G)
	expected := types.ZeroPoint()
	if result2.X.Cmp(expected.X) != 0 || result2.Y.Cmp(expected.Y) != 0 {
		t.Errorf("0*G != ZeroPoint, got (%x, %x)", result2.X, result2.Y)
	}

	// 2 * G = G + G
	G2 := AddPoints(G, G)
	result3 := ScalarMult(big.NewInt(2), G)
	if result3.X.Cmp(G2.X) != 0 || result3.Y.Cmp(G2.Y) != 0 {
		t.Errorf("2*G != G+G, got (%x, %x), expected (%x, %x)", result3.X, result3.Y, G2.X, G2.Y)
	}
}

func TestPubKeyFromPrivate(t *testing.T) {
	// priv key 1 -> pub key = G
	pub := PubKeyFromPrivate(big.NewInt(1))
	if pub.X.Cmp(Gx) != 0 || pub.Y.Cmp(Gy) != 0 {
		t.Errorf("private 1 gives wrong pubkey, got (%x, %x)", pub.X, pub.Y)
	}

	// priv key 2 -> pub key = 2*G
	G2 := AddPoints(&types.Point{X: Gx, Y: Gy}, &types.Point{X: Gx, Y: Gy})
	pub2 := PubKeyFromPrivate(big.NewInt(2))
	if pub2.X.Cmp(G2.X) != 0 || pub2.Y.Cmp(G2.Y) != 0 {
		t.Errorf("private 2 gives wrong pubkey, got (%x, %x)", pub2.X, pub2.Y)
	}
}

func TestParsePublicKey(t *testing.T) {
	compressed := "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	pub, err := ParsePubKey(compressed)
	if err != nil {
		t.Fatalf("ParsePublicKey failed: %v", err)
	}
	if pub.X.Cmp(Gx) != 0 || pub.Y.Cmp(Gy) != 0 {
		t.Errorf("parsed pubkey does not match G, got (%x, %x)", pub.X, pub.Y)
	}

	uncompressed := "0479be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	pub2, err := ParsePubKey(uncompressed)
	if err != nil {
		t.Fatalf("ParsePublicKey failed: %v", err)
	}
	if pub2.X.Cmp(Gx) != 0 || pub2.Y.Cmp(Gy) != 0 {
		t.Errorf("parsed pubkey does not match G, got (%x, %x)", pub2.X, pub2.Y)
	}

	_, err = ParsePubKey("invalid")
	if err == nil {
		t.Error("expected err invalid key")
	}
}
