package simpletcp

import (
	"context"
	"errors"
	"net"

	"golang.org/x/sync/errgroup"
	"gopkg.in/dailymuse/gzap.v1"
)

var (
	ErrNilConfig      = errors.New("configuration is nil")
	ErrNotImplemented = errors.New("not implemented")
	ErrNilQuotes      = errors.New("quotes service is nil")
)

type configurer interface {
	GetBindAddress() string
}

type quoterer interface {
	GetQuote(context.Context) (string, error)
}

type tcpserver struct {
	quotes      quoterer
	bindAddress string

	port int
}

func (s *tcpserver) Startup(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Listen for incoming connections.
	l, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		return err
	}
	s.port = l.Addr().(*net.TCPAddr).Port

	g.Go(func() error {
		<-ctx.Done()
		defer l.Close()
		return nil
	})

	g.Go(func() error {
		for {
			conn, err1 := l.Accept()
			if err1 != nil {
				l.Close()
				select {
				case <-ctx.Done():
					return nil
				default:
					return err
				}
			}
			s.handleConnection(ctx, g, conn)
		}
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s *tcpserver) GetPort() int {
	return s.port
}

func (s *tcpserver) handleConnection(ctx context.Context, g *errgroup.Group, conn net.Conn) {
	g.Go(func() error {

		if s.quotes != nil {

			quote, err := s.quotes.GetQuote(ctx)
			if err != nil {
				gzap.Logger.Error("error get quote", gzap.Error(err))
			} else {
				_, err = conn.Write([]byte(quote))
				if err != nil {
					gzap.Logger.Error("error write to client", gzap.Error(err))
				}
			}
		}

		err := conn.Close()
		if err != nil {
			gzap.Logger.Error("error close connection", gzap.Error(err))
			return nil
		}
		return nil
	})
}

func NewSimpleTCPServer(
	cfg configurer,
	quotes quoterer,
) (*tcpserver, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if quotes == nil {
		return nil, ErrNilQuotes
	}

	s := &tcpserver{
		bindAddress: cfg.GetBindAddress(),
		quotes:      quotes,
	}

	return s, nil
}
