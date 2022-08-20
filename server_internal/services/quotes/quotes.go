package quotes

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
	ErrZeroWorkersCount      = errors.New("zero workers count")
	ErrEmptyFilepath         = errors.New("empty quotes filepath")
	ErrNotImplemented        = errors.New("not implemented")
	ErrQuotesFileIsNotExists = errors.New("quotes file is not exists")
	ErrEmptyDatabase         = errors.New("empty database")
)

type configurer interface {
	GetQuotesFilepath() string
	GetProcessingTime() time.Duration
	GetWorkersCount() uint
}

type quotesService struct {
	quotesFilepath string
	processingTime time.Duration
	workersCount   uint

	quotes *models.Quotes

	mu deadlock.RWMutex
}

func (qs *quotesService) Startup(ctx context.Context) (err error) {
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

func (qs *quotesService) GetQuote() (string, error) {
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
