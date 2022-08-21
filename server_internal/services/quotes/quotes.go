package quotes

import (
	"context"
	"errors"

	"github.com/sasha-s/go-deadlock"
	"github.com/svkior/powgwey/server_internal/models"
	"golang.org/x/sync/errgroup"
	"gopkg.in/dailymuse/gzap.v1"
)

var (
	ErrNilConfig        = errors.New("configuration is nil")
	ErrZeroWorkersCount = errors.New("zero workers count")
	ErrNotImplemented   = errors.New("not implemented")
	ErrNotIninializated = errors.New("service is not running")
	ErrStorageIsNil     = errors.New("storage is nil")
	ErrShutdown         = errors.New("service is shutdowning")
)

type configurer interface {
	GetWorkersCount() uint
}

type quotesStorager interface {
	GetQuote(ctx context.Context) (string, error)
}

type quotesService struct {
	workersCount uint
	storage      quotesStorager

	workerChannel chan chan *models.QuotesWork
	input         chan *models.QuotesWork
	mu            deadlock.RWMutex
	cancel        func()
}

func (qs *quotesService) Startup(ctx context.Context) (err error) {
	ctx, qs.cancel = context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return qs.initWorkerPool(ctx, g)
	})

	g.Go(func() error {
		<-ctx.Done()
		return nil
	})
	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (qs *quotesService) Shutdown(_ context.Context) (err error) {
	if qs.cancel == nil {
		return ErrNotIninializated
	}

	qs.cancel()

	return nil
}

func (qs *quotesService) GetQuote(ctx context.Context) (string, error) {
	if qs.storage == nil {
		return "", ErrNotIninializated
	}

	if qs.input == nil {
		return "", ErrNotIninializated
	}

	finish := make(chan struct{})
	job := &models.QuotesWork{
		Finish: finish,
	}

	select {
	case <-ctx.Done():
		return "", ErrShutdown
	case qs.input <- job:
		select {
		case <-ctx.Done():
			return "", ErrShutdown
		case <-finish:
			if job.Error != nil {
				return "", job.Error
			}
			return job.Quote, nil
		}
	}
}

func (qs *quotesService) initWorkerPool(ctx context.Context, g *errgroup.Group) error {
	for i := uint(0); i < qs.workersCount; i++ {
		gzap.Logger.Info("starting worker",
			gzap.Uint("ID", i),
		)

		w, err := NewQuotesWorker(
			qs.storage,
			qs.workerChannel,
		)
		if err != nil {
			gzap.Logger.Error("error starting worker",
				gzap.Error(err),
				gzap.Uint("ID", i),
			)
			return err
		}

		g.Go(func() error {
			err = w.Startup(ctx)
			if err != nil {
				return err
			}

			return nil
		})
	}

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case work := <-qs.input:
				select {
				case <-ctx.Done():
					return nil
				case currentWorker := <-qs.workerChannel:
					select {
					case <-ctx.Done():
						return nil
					case currentWorker <- work:
						continue
					}
				}
			}
		}
	})

	return nil
}

func NewQuotesService(
	ctx context.Context,
	cfg configurer,
	storage quotesStorager,
) (*quotesService, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if storage == nil {
		return nil, ErrStorageIsNil
	}

	s := &quotesService{
		input:         make(chan *models.QuotesWork),
		storage:       storage,
		workersCount:  cfg.GetWorkersCount(),
		workerChannel: make(chan chan *models.QuotesWork),
	}

	if s.workersCount < 1 {
		return nil, ErrZeroWorkersCount
	}

	return s, nil
}
