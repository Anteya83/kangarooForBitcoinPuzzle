package config

type Config struct {
	PublicKeyHex      string
	StartRangeHex     string
	EndRangeHex       string
	NumWorkers        int
	MaxJump           int64 // 2^MaxJump
	TameStepsLimit    int64
	DistinguishedBits int
	NumShards         int
	SaveInterval      int64
}

// have to change PublicKeyHex and Key Range: --StartRangeHex-EndRangeHex
func DefaultConfig() Config {
	return Config{
		PublicKeyHex:      "033c4a45cbd643ff97d77f41ea37e843648d50fd894b864b0d52febc62f6454f7c",
		StartRangeHex:     "80000",
		EndRangeHex:       "fffff",
		NumWorkers:        16,
		MaxJump:           0,
		TameStepsLimit:    1000000000000000000,
		DistinguishedBits: 8,
		NumShards:         64,
		SaveInterval:      10000,
	}
}
