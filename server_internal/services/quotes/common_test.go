package quotes_test

import (
	"context"
	"time"
)

const (
	TimeoutOnStartup = 300 * time.Millisecond
)

type configMock struct {
	numberOfWorkers uint
}

func (c *configMock) GetWorkersCount() uint {
	return c.numberOfWorkers
}

type storageMock struct {
	quote string
	err   error
}

func (s *storageMock) GetQuote(ctx context.Context) (string, error) {
	return s.quote, s.err
}
