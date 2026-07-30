package encoding

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func PrivateKeyToWIF(privateKeyHex string) (string, error) {
	keyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}
	versionedKey := append([]byte{0x80}, keyBytes...)
	checksum := calculateChecksum(versionedKey)
	fullKey := append(versionedKey, checksum...)
	return base58.Encode(fullKey), nil
}

func GetPublicKey(privateKeyHex string) ([]byte, error) {
	keyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(keyBytes)
	if privateKey == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return privateKey.PubKey().SerializeUncompressed(), nil
}

func CompressPublicKey(uncompressedKey []byte) []byte {
	if len(uncompressedKey) != 65 || uncompressedKey[0] != 0x04 {
		return uncompressedKey
	}
	x := uncompressedKey[1:33]
	y := uncompressedKey[33:]
	prefix := byte(0x02)
	if y[len(y)-1]%2 != 0 {
		prefix = 0x03
	}
	return append([]byte{prefix}, x...)
}

func PublicKeyToAddress(publicKey []byte) (string, error) {
	sha256Hash := sha256.Sum256(publicKey)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		return "", err
	}
	ripemd160Hash := ripemd160Hasher.Sum(nil)
	versionedPayload := append([]byte{0x00}, ripemd160Hash...)
	checksum := calculateChecksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	return base58.Encode(fullPayload), nil
}

func PublicKeyToP2SHAddress(publicKey []byte) (string, error) {
	compressedPubKey := CompressPublicKey(publicKey)
	sha1 := sha256.Sum256(compressedPubKey)
	ripemd1 := ripemd160.New()
	ripemd1.Write(sha1[:])
	ripemd1Hash := ripemd1.Sum(nil)

	withPrefix := append([]byte{0x00, 0x14}, ripemd1Hash...)
	sha2 := sha256.Sum256(withPrefix)
	ripemd2 := ripemd160.New()
	ripemd2.Write(sha2[:])
	ripemd2Hash := ripemd2.Sum(nil)

	withPrefix2 := append([]byte{0x05}, ripemd2Hash...)
	checksum := calculateChecksum(withPrefix2)
	finalPayload := append(withPrefix2, checksum...)
	return base58.Encode(finalPayload), nil
}

func calculateChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}
