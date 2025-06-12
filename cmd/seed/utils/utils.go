package utils

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
)

func GenerateNumberBetween(min int, max int) int {
	cryptoRand, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		// Fallback to math/rand if crypto/rand fails
		return min + rand.Intn(max-min+1)
	}
	return min + int(cryptoRand.Int64())
}
