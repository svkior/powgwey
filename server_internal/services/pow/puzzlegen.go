package pow

import (
	"crypto/rand"
	"math/big"

	"github.com/svkior/powgwey/server_internal/models"
)

func randUint64(ceil uint64) (uint64, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(ceil)))
	if err != nil {
		return 0, err
	}
	return val.Uint64(), nil
}

func generatePuzzle(
	algo func(uint64, float64) uint64,
	n, k int64,
) (uint64, *models.Puzzle, error) {
	max := (uint64(1) << n) - 1
	maxF := float64(max)

	x0, err := randUint64(max)
	if err != nil {
		return 0, nil, err
	}

	seq := make([]uint64, k+1)
	seq[k] = x0

	xk := x0
	for idx := uint64(1); idx <= uint64(k); idx++ {
		xk = algo(xk, maxF) ^ idx
		seq[uint64(k)-idx] = xk
	}
	checkSum := checksum(seq)

	return x0, &models.Puzzle{
		Xk:       xk,
		K:        k,
		N:        n,
		Checksum: checkSum,
	}, nil
}
