package server

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

// requestURL is a data structure which represents a shortening request
type requestURL struct {
	URL string `json:"long_url"`
}

// encodeLongURL generates an ID for a given long URL.
func (u *requestURL) encodeLongURL() (string, error) {
	toSHA256 := toSHA256(u.URL)
	toB64Numeric := fmt.Sprint(new(big.Int).SetBytes(toSHA256).Uint64())
	finalString, err := toBase58([]byte(toB64Numeric))
	if err != nil {
		return "", err
	}
	return finalString[:8], nil
}

// toSHA256 converts a string to SHA256 bytes
func toSHA256(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// toBase58 converts a slice of bytes to base58 encoding to make it look nice.
func toBase58(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
