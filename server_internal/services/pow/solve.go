package pow

import (
	"errors"

	"github.com/svkior/powgwey/server_internal/models"
)

var (
	ErrSolutionNotFound = errors.New("solution not found")
)

type solver struct {
	invTable map[uint64][]uint64

	*models.Puzzle
}

func Solution(
	p *models.Puzzle,
	algo func(uint64, float64) uint64,
) (*models.Solve, error) {

	s := &solver{
		invTable: generateInversionTable(p.N, algo),
		Puzzle:   p,
	}

	solution, hasSolution := s.findSolution(p.Xk, []uint64{})
	if !hasSolution {
		return nil, ErrSolutionNotFound
	}

	return &models.Solve{
		Y0: solution,
	}, nil
}

func (s *solver) findSolution(
	curVal uint64,
	seq []uint64,
) (uint64, bool) {

	currDepth := s.K - int64(len(seq))
	seq = append(seq, curVal)

	if currDepth == 0 {
		checksumToCheck := checksum(seq)

		if checksumToCheck == s.Checksum {
			return curVal, true
		}

		return 0, false
	}

	cvcd := curVal ^ uint64(currDepth)
	leafs, hasLeafs := s.invTable[cvcd]
	if !hasLeafs {
		return 0, false
	}

	for _, leaf := range leafs {
		value, ok := s.findSolution(leaf, seq)
		if ok {
			return value, true
		}
	}

	return 0, false

}

func generateInversionTable(n int64, algo func(uint64, float64) uint64) map[uint64][]uint64 {
	max := (uint64(1) << n) - 1
	maxF := float64(max)
	inversionTable := make(map[uint64][]uint64, max)

	var res uint64
	for idx := uint64(0); idx <= max; idx++ {
		res = algo(idx, maxF)

		row, hasRow := inversionTable[res]
		if !hasRow {
			row = []uint64{}
		}

		inversionTable[res] = append(row, idx)
	}
	return inversionTable
}
