package quotes

import (
	"context"
	"errors"

	"github.com/svkior/powgwey/server_internal/models"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNilWorkerChannel = errors.New("nil worker channel")
)

type worker struct {
	storage       quotesStorager
	workerChannel chan chan *models.QuotesWork
	channel       chan *models.QuotesWork
	cancel        func()
}

func (w *worker) Startup(ctx context.Context) error {
	ctx, w.cancel = context.WithCancel(ctx)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case w.workerChannel <- w.channel:
				select {
				case <-ctx.Done():
					return nil
				case job := <-w.channel:
					err := w.DoWork(ctx, job)
					if err != nil {
						return err
					}
				}
			}
		}
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) DoWork(ctx context.Context, job *models.QuotesWork) error {
	defer func() { job.Finish <- struct{}{} }()

	if w.storage == nil {
		return ErrStorageIsNil
	}

	job.Quote, job.Error = w.storage.GetQuote(ctx)

	return nil
}

func NewQuotesWorker(
	storage quotesStorager,
	workerChannel chan chan *models.QuotesWork,
) (*worker, error) {

	if storage == nil {
		return nil, ErrStorageIsNil
	}

	if workerChannel == nil {
		return nil, ErrNilWorkerChannel
	}

	w := &worker{
		storage:       storage,
		workerChannel: workerChannel,
		channel:       make(chan *models.QuotesWork),
	}

	return w, nil
}
