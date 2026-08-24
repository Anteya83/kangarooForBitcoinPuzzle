package cli

import (
	"flag"
	"kangarooForBitcoinPuzzle/internal/config"
)

func ParseFlags() config.Config {
	cfg := config.DefaultConfig()

	var (
		pubKeyHex = flag.String("pubkey", cfg.PublicKeyHex, "Public key in hex (compressed or uncompressed)")
		startHex  = flag.String("start", cfg.StartRangeHex, "Start of range in hex")
		endHex    = flag.String("end", cfg.EndRangeHex, "End of range in hex")
		workers   = flag.Int("workers", cfg.NumWorkers, "Number of worker goroutines")
		dpBits    = flag.Int("dpbits", cfg.DistinguishedBits, "Distinguished bits (number of low bits to check)")
		maxJump   = flag.Int64("maxjump", cfg.MaxJump, "Max jump size (0 = auto)")
		tameLimit = flag.Int64("tamelimit", cfg.TameStepsLimit, "Tame steps limit (0 = no limit, but practically 1e18)")
	)
	flag.Parse()

	cfg.PublicKeyHex = *pubKeyHex
	cfg.StartRangeHex = *startHex
	cfg.EndRangeHex = *endHex
	cfg.NumWorkers = *workers
	cfg.DistinguishedBits = *dpBits
	cfg.MaxJump = *maxJump
	cfg.TameStepsLimit = *tameLimit

	return cfg
}
