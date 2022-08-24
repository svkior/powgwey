package pow

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/svkior/powgwey/server_internal/models"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNilConfig      = errors.New("nil config")
	ErrNilMetrics     = errors.New("nil metrics")
	ErrNilConnection  = errors.New("nil connection")
	ErrTimeout        = errors.New("io timeout")
	ErrGetRequest     = errors.New("wrong request")
	ErrGeneratePuzzle = errors.New("error generating puzzle")
	ErrGetSolution    = errors.New("error get solution")
	ErrWrongAnswer    = errors.New("client got wrong  answer")
)

type configurer interface {
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}

type metricker interface {
	GetDifficulties() uint
}

type pow struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	metric       metricker
}

//nolint:gocyclo  // Big saga
func (m *pow) Validate(ctx context.Context, conn net.Conn) error {
	if conn == nil {
		return ErrNilConnection
	}
	reply := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(m.readTimeout))
	bytes_readed, err := conn.Read(reply)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}
		if os.IsTimeout(err) {
			return ErrTimeout
		}
		return err
	}
	conn.SetReadDeadline(time.Now().Add(10 * m.readTimeout))
	conn.SetWriteDeadline(time.Now().Add(10 * m.readTimeout))

	req := models.Request{}
	err = req.UnmarshalBinary(reply[0:bytes_readed])
	if err != nil {
		return ErrGetRequest
	}

	x0, puzzle, err := m.generatePuzzle()
	if err != nil {
		return ErrGeneratePuzzle
	}

	marshalLen := puzzle.MarshalTo(reply)
	if err != nil {
		return ErrGeneratePuzzle
	}
	conn.SetReadDeadline(time.Now().Add(100 * m.readTimeout))
	conn.SetWriteDeadline(time.Now().Add(m.writeTimeout))
	_, err = conn.Write(reply[0:marshalLen])
	if err != nil {
		return ErrTimeout
	}

	conn.SetWriteDeadline(time.Now().Add(100 * m.writeTimeout))
	conn.SetReadDeadline(time.Now().Add(100 * m.readTimeout))
	bytes_readed, err = conn.Read(reply)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}
		if os.IsTimeout(err) {
			return ErrTimeout
		}
		return err
	}

	answer := models.Solve{}
	err = answer.UnmarshalBinary(reply[0:bytes_readed])
	if err != nil {
		return ErrGetSolution
	}

	if x0 != answer.Y0 {
		return ErrWrongAnswer
	}

	ok := &models.Result{Ok: true}
	marshalLen = ok.MarshalTo(reply)
	_, err = conn.Write(reply[0:marshalLen])
	if err != nil {
		return ErrTimeout
	}

	return nil
}

func (m *pow) generatePuzzle() (uint64, *models.Puzzle, error) {

	difficulty := m.metric.GetDifficulties()

	log.Printf("difficulty %v", difficulty)
	//FIXME: Choose difficulty of algorythn  from  difficulty
	x0, puzzle, err := generatePuzzle(Algosin, 5, 5)
	if err != nil {
		return 0, nil, err
	}

	return x0, puzzle, nil
}

func NewPoWMiddleware(
	cfg configurer,
	metrics metricker,
) (*pow, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if metrics == nil {
		return nil, ErrNilMetrics
	}

	p := &pow{
		readTimeout:  cfg.GetReadTimeout(),
		writeTimeout: cfg.GetWriteTimeout(),
		metric:       metrics,
	}
	return p, nil
}
