package config

import (
	"fmt"
	"math/big"

	"kangarooForBitcoinPuzzle/internal/crypto"
)

func Validate(cfg Config) error {
	if cfg.NumWorkers <= 0 {
		return fmt.Errorf("NumWorkers must be > 0, got %d", cfg.NumWorkers)
	}

	if cfg.NumShards <= 0 {
		return fmt.Errorf("NumShards must be > 0, got %d", cfg.NumShards)
	}

	if cfg.DistinguishedBits < 1 {
		return fmt.Errorf("DistinguishedBits must be >= 1, got %d", cfg.DistinguishedBits)
	}

	a := new(big.Int)
	b := new(big.Int)
	if _, ok := a.SetString(cfg.StartRangeHex, 16); !ok {
		return fmt.Errorf("invalid StartRangeHex: %s", cfg.StartRangeHex)
	}
	if _, ok := b.SetString(cfg.EndRangeHex, 16); !ok {
		return fmt.Errorf("invalid EndRangeHex: %s", cfg.EndRangeHex)
	}
	if a.Cmp(b) >= 0 {
		return fmt.Errorf("StartRangeHex (%s) must be less than EndRangeHex (%s)", cfg.StartRangeHex, cfg.EndRangeHex)
	}

	if cfg.SaveInterval <= 0 {
		return fmt.Errorf("SaveInterval must be > 0, got %d", cfg.SaveInterval)
	}

	if cfg.MaxJump < 0 {
		return fmt.Errorf("MaxJump cannot be negative, got %d", cfg.MaxJump)
	}

	if cfg.TameStepsLimit <= 0 {
		return fmt.Errorf("TameStepsLimit must be > 0, got %d", cfg.TameStepsLimit)
	}

	if _, err := crypto.ParsePubKey(cfg.PublicKeyHex); err != nil {
		return fmt.Errorf("invalid PublicKeyHex: %v", err)
	}

	return nil
}
