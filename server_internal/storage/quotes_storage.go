package storage

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/svkior/powgwey/server_internal/app"
	"github.com/svkior/powgwey/server_internal/models"

	"github.com/mailru/easyjson"
	"github.com/sasha-s/go-deadlock"
	"golang.org/x/sync/errgroup"
	"gopkg.in/dailymuse/gzap.v1"
)

var (
	ErrNilConfig             = errors.New("configuration is nil")
	ErrEmptyFilepath         = errors.New("empty quotes filepath")
	ErrQuotesFileIsNotExists = errors.New("quotes file is not exists")
	ErrStorageIsNotStarted   = errors.New("can't stop storage is not started")
	ErrEmptyDatabase         = errors.New("empty database")
	ErrNotImplemented        = errors.New("not implemented")
	ErrShuttingDown          = errors.New("shutting down storage")
)

type configurer interface {
	GetQuotesFilepath() string
	GetProcessingTime() time.Duration
}

type quotesStorage struct {
	quotesFilepath string
	processingTime time.Duration

	quotes *models.Quotes

	mu deadlock.RWMutex
}

func (qs *quotesStorage) Startup(ctx context.Context) (err error) {

	gzap.Logger.Info("Starting quotes storage",
		gzap.String("quotes from", qs.quotesFilepath),
		gzap.Duration("processing time", qs.processingTime),
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return qs.loadQuotes(ctx)
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

func (qs *quotesStorage) GetQuote(ctx context.Context) (string, error) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	if len(*qs.quotes) == 0 {
		return "", ErrEmptyDatabase
	}

	numberOfQuotes := int64(len(*qs.quotes)) - 1

	result, _ := rand.Int(rand.Reader, big.NewInt(numberOfQuotes))

	index := result.Int64()

	quote := (*qs.quotes)[index]

	select {
	case <-ctx.Done():
		return "", ErrShuttingDown
	case <-time.After(qs.processingTime):
	}

	return quote.Quote, nil
}

func (qs *quotesStorage) loadQuotes(ctx context.Context) error {
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

func NewQuotesStorage(
	ctx context.Context,
	cfg configurer,
) (*quotesStorage, error) {

	if cfg == nil {
		return nil, ErrNilConfig
	}

	s := &quotesStorage{
		quotesFilepath: cfg.GetQuotesFilepath(),
		processingTime: cfg.GetProcessingTime(),
	}

	if s.quotesFilepath == "" {
		return nil, ErrEmptyFilepath
	}

	if !fileExists(s.quotesFilepath) {
		return nil, ErrQuotesFileIsNotExists
	}

	return s, nil
}
