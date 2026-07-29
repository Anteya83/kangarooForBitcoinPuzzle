package main

import (
	"fmt"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"math/big"
)

func main() {
	fmt.Println("--Kangaroo start--")

	priv := big.NewInt(10000)
	pub := crypto.PubKeyFromPrivate(priv)
	fmt.Printf("Privkey: %v\n", priv)
	fmt.Printf("Pub X: %x\n", pub.X)
	fmt.Printf("Pub Y: %x\n", pub.Y)

	compressed := "024ee2be2d4e9f92d2f5a4a03058617dc45befe22938feed5b7a6b7282dd74cbdd"
	pub2, err := crypto.ParsePubKey(compressed)
	if err != nil {
		fmt.Println("parsing crash:", err)
	} else {
		fmt.Printf("Pub X: %x\n", pub2.X)
		fmt.Printf("Pub Y: %x\n", pub2.Y)
	}
}
