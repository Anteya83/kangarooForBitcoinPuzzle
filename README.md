# Kangaroo Solver for Bitcoin Private Keys

A concurrent implementation of **Pollard's Kangaroo (Lambda) algorithm** in Go, designed to find a Bitcoin private key for a given public key when the private key is known to lie within a specified range. This tool is useful for cryptanalysis challenges, Bitcoin puzzles, and educational purposes.

---

## Features

- **Efficient Kangaroo algorithm** with distinguished points.
- **Automatic jump size optimization** based on the range width.
- **Multi‑worker concurrency** (goroutines) for faster search.
- **State persistence**: saves and resumes progress automatically.
- **Full Bitcoin key utilities**:
  - Private key to WIF (compressed)
  - Public key extraction and compression
  - P2PKH and P2SH address generation
- **CLI flags** for flexible configuration.
- **Clean package structure** for easy extension.

---

## Project Structure

```
kangarooForBitcoinPuzzle/
├── main.go                  # Entry point
├── go.mod
├── go.sum
├── .gitignore
└── internal/
    ├── types/                # Point on elliptic curve
    ├── crypto/               # EC arithmetic, key parsing
    ├── shardedmap/           # Concurrent map for distinguished points
    ├── state/                # Save/load persistent state
    ├── config/               # Configuration struct and defaults
    ├── encoding/             # WIF, addresses
    ├── kangaroo/             # Main solver logic
    └── cli/                  # CLI flags parsing
```
---
### Prerequisites

- Go 1.26 or later
---

### Build and Run

Clone the repository and build:

```bash
git clone https://github.com/Anteya83/kangarooForBitcoinPuzzle.git
cd kangarooForBitcoinPuzzle
go mod tidy
```

Run the solver with default settings (which are set for a specific puzzle range):

```bash
go run main.go
```

Or build and run the binary:

```bash
go build -o kangaroo main.go
./kangaroo
```

### Running Tests

To verify the algorithm works correctly:

```bash
go test ./internal/kangaroo/ -v
```

This runs a small-range test that should complete in under 2 seconds.

---
## CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-pubkey` | Public key in hex (compressed or uncompressed) | (from config) |
| `-start` | Start of the search range (hexadecimal) | (from config) |
| `-end` | End of the search range (hexadecimal) | (from config) |
| `-workers` | Number of concurrent goroutines (wild kangaroos) | `16` |
| `-dpbits` | Distinguished bits (low bits of X to check) | `12` |
| `-maxjump` | Max jump size (`0` = auto‑calculate) | `0` (auto) |
| `-tamelimit` | Maximum steps for the tame kangaroo before giving up | `1e18` (effectively unlimited) |
To see all flags:

```bash
./kangaroo -help
````
---
## Configuration

All configurable parameters are defined in `internal/config/config.go`. You can either modify the `DefaultConfig()` function .

### Key Parameters

| Parameter            | Description                                                  | Default Example               |
|----------------------|--------------------------------------------------------------|-------------------------------|
| `PublicKeyHex`       | The target public key (compressed or uncompressed, hex).     | `"024ee2be2d4e9f..."`         |
| `StartRangeHex`      | Lower bound of the search range (hexadecimal).               | `"28299999999000000"`         |
| `EndRangeHex`        | Upper bound of the search range (hexadecimal).               | `"2840000000000ffff"`         |
| `NumWorkers`         | Number of concurrent goroutines (wild kangaroos).            | `16`                          |
| `DistinguishedBits`  | Number of low bits that define a distinguished point.        | `12`                          |
| `TameStepsLimit`     | Maximum steps for the tame kangaroo before giving up.        | `1e18` (effectively unlimited)|
| `SaveInterval`       | How often (in new distinguished points) to save state.       | `10000`                       |
| `NumShards`          | Number of shards for the concurrent map.                     | `64`                          |

### How to Change

**Option A: Edit `config.go` directly**

Modify the `DefaultConfig()` function in `internal/config/config.go` to set your own values.

**Option B: Use CLI flags (recommended)**

```bash
./kangaroo -pubkey 024ee... -start 80000 -end fffff -workers 8 -dpbits 4 -tamelimit 1000000
````
## How It Works

1. The **tame kangaroo** starts from the upper bound `b * G` and makes pseudorandom jumps, storing only **distinguished points** (points whose X coordinate has a certain number of low bits zero).
2. The **wild kangaroo** starts from the target public key and makes pseudorandom upward jumps.
3. When a wild kangaroo lands on a distinguished point already visited by the tame kangaroo, a collision occurs.
4. The private key is recovered as `b - tame_distance - wild_distance`.
5. The solver automatically saves its progress to `state_<hash>.bin`, so you can stop and resume anytime.

---

## Output Example

```
Pub key: 033c4a45cbd643ff97d77f41ea37e843648d50fd894b864b0d52febc62f6454f7c
Range: 80000 - fffff
Auto MaxJump = 10 (width bits=19)
Starting Kangaroo search (Distinguished bits: 12)
Workers: 16
Range: 80000 - fffff (HEX)
State file: state_f2066a1acc0c4735.bin
Starting fresh search...
Populating distinguished points...
Runtime: 1.221665s
Collision time: 1.220648584s
Tame steps: 9
Wild steps: 31474
Collisions found: 1
PrivKey HEX: 00000000000000000000000000000000000000000000000000000000000d2c55
WIF(compressed): 5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEtMbVN6VkAe
Addr(P2PKH): 1HsMJxNiV7TLxmoF6uJNkydxPFDog4NQum
Addr(P2SH): 3E6zPtGL2sRfrQijjZt4hD7vs1HkW2vB7m
Time: 1.221660625s
```

---

## Important Notes

- This is a **CPU‑only** implementation. For ranges larger than ~2⁶⁰, consider using GPU‑optimized versions (e.g., Jean‑Luc PONS' Kangaroo).
- The algorithm is **probabilistic** – it may fail if the tame kangaroo doesn't accumulate enough distinguished points. Increase `TameStepsLimit` if needed.
- State files are tied to the exact public key and range. If you change the range, a new state file will be created.

---

## License

This project is open‑source and available under the [MIT License](LICENSE).

---


## Contributing

Feel free to submit issues and pull requests. For major changes, please open an issue first to discuss what you would like to change.

---
## Donations
If you find this tool useful or if you simply liked , donations are welcome!

Bitcoin (BTC): bc1q2xmq60x9qg5xfs5t8aavghj4jfzqknuq5plz6x

---
## About me

I've been interested in Bitcoin for a long time. 
I have ready‑made programs for brute‑forcing and seed phrase recovery when some known data is available. 
I can write a custom program tailored to your specific request.

