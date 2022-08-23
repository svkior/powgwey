package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
)

// Simple F(x) = F(x') function
func Algosin(x uint64, max float64) uint64 {
	result := max * math.Abs(math.Sin(float64(x)))
	return uint64(result)
}

func checksum(src []uint64) string {
	hexString := make([]byte, 0, len(src)*16)

	for _, v := range src {
		hexString = append(hexString, []byte(fmt.Sprintf("%016x", v))...)
	}

	hexHash := sha256.Sum256(hexString)
	hash := hex.EncodeToString(hexHash[:])

	return hash
}
