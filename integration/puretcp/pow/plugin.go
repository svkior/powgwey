package pow

import (
	"errors"
	"net"

	"puretcp/models"
)

var (
	ErrNotSolved = errors.New("puzzle is not solved")
)

type solverPlugin struct {
}

func (p *solverPlugin) Solve(conn net.Conn) error {
	reply := make([]byte, 8192)
	req := models.Request{}
	bytes_marshalled := req.MarshalTo(reply)
	//SetWriteTimeout
	_, err := conn.Write(reply[0:bytes_marshalled])
	if err != nil {
		return err
	}
	bytes_readed, err := conn.Read(reply)
	if err != nil {
		return err
	}

	puzzle := &models.Puzzle{}
	_, err = puzzle.Unmarshal(reply[0:bytes_readed])
	if err != nil {
		return err
	}

	solution, err := Solution(puzzle, Algosin)
	bytes_marshalled = solution.MarshalTo(reply)
	_, err = conn.Write(reply[0:bytes_marshalled])
	if err != nil {
		return err
	}

	bytes_readed, err = conn.Read(reply)
	if err != nil {
		return err
	}

	ok := &models.Result{}
	_, err = ok.Unmarshal(reply[0:bytes_readed])
	if err != nil {
		return err
	}

	if !ok.Ok {
		return ErrNotSolved
	}
	return nil
}

func NewSolverPlugin() (*solverPlugin, error) {
	p := &solverPlugin{}

	return p, nil
}
