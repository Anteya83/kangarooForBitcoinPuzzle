package encoding

import (
	"testing"
)

func TestPrivateKeyToWIF(t *testing.T) {
	tests := []struct {
		name     string
		privHex  string
		expected string
	}{
		{
			name:     "private key 1 (uncompressed WIF)",
			privHex:  "0000000000000000000000000000000000000000000000000000000000000001",
			expected: "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAnchuDf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wif, err := PrivateKeyToWIF(tt.privHex)
			if err != nil {
				t.Fatal(err)
			}
			if wif != tt.expected {
				t.Errorf("PrivateKeyToWIF(%q) = %q, want %q", tt.privHex, wif, tt.expected)
			}
		})
	}
}

func TestPublicKeyToAddress(t *testing.T) {
	privHex := "0000000000000000000000000000000000000000000000000000000000000001"
	pub, err := GetPublicKey(privHex)
	if err != nil {
		t.Fatal(err)
	}
	compressed := CompressPublicKey(pub)
	addr, err := PublicKeyToAddress(compressed)
	if err != nil {
		t.Fatal(err)
	}
	expected := "1BgGZ9tcN4rm9KBzDn7KprQz87SZ26SAMH"
	if addr != expected {
		t.Errorf("Address = %q, want %q", addr, expected)
	}
}

func TestPublicKeyToP2SHAddress(t *testing.T) {
	privHex := "0000000000000000000000000000000000000000000000000000000000000001"
	pub, err := GetPublicKey(privHex)
	if err != nil {
		t.Fatal(err)
	}
	compressed := CompressPublicKey(pub)
	addr, err := PublicKeyToP2SHAddress(compressed)
	if err != nil {
		t.Fatal(err)
	}
	if len(addr) < 30 || addr[0] != '3' {
		t.Errorf("Invalid P2SH address: %q", addr)
	}
}
