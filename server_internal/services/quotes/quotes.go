package quotes

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"os"

	"github.com/svkior/powgwey/server_internal/app"
	"github.com/svkior/powgwey/server_internal/models"

	"github.com/mailru/easyjson"
	"github.com/sasha-s/go-deadlock"
	"golang.org/x/sync/errgroup"
	"gopkg.in/dailymuse/gzap.v1"
)

var (
	ErrNilConfig        = errors.New("configuration is nil")
	ErrZeroWorkersCount = errors.New("zero workers count")
	ErrNotImplemented   = errors.New("not implemented")
)

type configurer interface {
	GetWorkersCount() uint
}

type quotesStorage interface {
	GetQuote(ctx context.Context) (string, error)
}

type quotesService struct {
	workersCount uint
	storage      quotesStorage

	mu     deadlock.RWMutex
	cancel func()
}

func (qs *quotesService) Startup(ctx context.Context) (err error) {
	ctx, qs.cancel = context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return qs.initWorkerPool(ctx)
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

func (qs *quotesService) 

func (qs *quotesService) GetQuote(ctx context.Context) (string, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	if len(*qs.quotes) == 0 {
		return "", ErrEmptyDatabase
	}

	numberOfQuotes := int64(len(*qs.quotes)) - 1

	result, _ := rand.Int(rand.Reader, big.NewInt(numberOfQuotes))

	index := result.Int64()

	quote := (*qs.quotes)[index]

	return quote.Quote, nil
}

func (qs *quotesService) loadQuotes(ctx context.Context) error {
	rawBytes, err := os.ReadFile(qs.quotesFilepath)
	if err != nil {
		gzap.Logger.Error("error reading quotes file",
			gzap.Error(err),
			gzap.String(app.FilePathTag, qs.quotesFilepath))
		return err
	}

	quotes := &models.Quotes{}
	err = easyjson.Unmarshal(rawBytes, quotes)
	if err != nil {
		gzap.Logger.Error("error unmarshal quotes file",
			gzap.Error(err),
			gzap.String(app.FilePathTag, qs.quotesFilepath))
		return err
	}

	qs.mu.Lock()
	defer qs.mu.Unlock()

	qs.quotes = quotes

	return nil
}

func NewQuotesService(
	ctx context.Context,
	cfg configurer,
) (*quotesService, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	s := &quotesService{
		quotesFilepath: cfg.GetQuotesFilepath(),
		processingTime: cfg.GetProcessingTime(),
		workersCount:   cfg.GetWorkersCount(),
	}

	if s.workersCount < 1 {
		return nil, ErrZeroWorkersCount
	}

	if len(s.quotesFilepath) == 0 {
		return nil, ErrEmptyFilepath
	}

	if !fileExists(s.quotesFilepath) {
		return nil, ErrQuotesFileIsNotExists
	}

	return s, nil
}
