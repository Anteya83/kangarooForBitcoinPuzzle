package main

import (
	"fmt"
	"kangarooForBitcoinPuzzle/internal/config"
	"kangarooForBitcoinPuzzle/internal/crypto"
	"kangarooForBitcoinPuzzle/internal/encoding"
	"kangarooForBitcoinPuzzle/internal/kangaroo"
	"math/big"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.DefaultConfig()

	fmt.Printf("Pub key: %s\n", cfg.PublicKeyHex)
	fmt.Printf("Range: %s - %s\n", cfg.StartRangeHex, cfg.EndRangeHex)

	pubKey, err := crypto.ParsePubKey(cfg.PublicKeyHex)
	if err != nil {
		panic(err)
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
		fmt.Printf("WIF(compressed): %s\n", wif)
		fmt.Printf("Addr(P2PKH): %s\n", addr)
		fmt.Printf("Addr(P2SH): %s\n", p2sh)
	} else {
		fmt.Println("Key not found")
	}
	fmt.Printf("Time: %v\n", elapsed)
}
