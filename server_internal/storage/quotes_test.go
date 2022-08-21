package storage_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/svkior/powgwey/server_internal/services/quotes"
	"github.com/svkior/powgwey/server_internal/storage"
	"golang.org/x/sync/errgroup"
)

type testConfig struct {
	processingTime time.Duration
	quotesFilepath string
}

func (c *testConfig) GetProcessingTime() time.Duration {
	return c.processingTime
}

func (c *testConfig) GetQuotesFilepath() string {
	return c.quotesFilepath
}

func TestSpect(t *testing.T) {
	ctx := context.TODO()
	Convey("Given quotes service with nil config", t, func() {
		_, err := storage.NewQuotesStorage(ctx, nil)
		Convey("The error should be ErrNilConfig", func() {
			So(err, ShouldResemble, quotes.ErrNilConfig)
		})
	})

	Convey("Given non empty config for quotes", t, func() {
		cfg := &testConfig{
			processingTime: 1,
			quotesFilepath: "../../data/quotes/movies.json",
		}
		Convey("We create new Quotes Service with zero processing time", func() {
			cfg.processingTime = 0
			_, err := storage.NewQuotesStorage(ctx, cfg)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("We create new Quotes with empty filepath", func() {
			cfg.quotesFilepath = ""
			_, err := storage.NewQuotesStorage(ctx, cfg)
			Convey("Error should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Error should be Empty quotes filepath", func() {
				So(err, ShouldEqual, storage.ErrEmptyFilepath)
			})
		})

		Convey("When we create new service with non existing filepath", func() {
			cfg.quotesFilepath = "non-existent-file"
			_, err := storage.NewQuotesStorage(ctx, cfg)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Error should be Empty quotes filepath", func() {
				So(err, ShouldEqual, storage.ErrQuotesFileIsNotExists)
			})
		})

		Convey("When we create new service with wrong json structure", func() {
			cfg.quotesFilepath = "../../data/quotes/wrong.json"
			qs, err := storage.NewQuotesStorage(ctx, cfg)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Service should NOT be nil", func() {
				So(qs, ShouldNotBeNil)
			})
			Convey("When we start this service", func() {
				ctx1, cancel := context.WithCancel(ctx)
				g, gCtx := errgroup.WithContext(ctx1)
				var err1 error

				g.Go(func() error {
					err1 = qs.Startup(gCtx)
					return err1
				})
				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return nil
					case <-time.After(1 * time.Second):
						cancel()
						return errors.New("timeout")
					}
				})
				err2 := g.Wait()
				Convey("should not be nil", func() {
					So(err2, ShouldNotBeNil)
				})
				Convey("should starts with parse", func() {
					So(fmt.Sprintf("%s", err2), ShouldStartWith, "parse")
				})
			})
		})

		Convey("When we create new Quotes Service with normal params", func() {
			qs, err := storage.NewQuotesStorage(ctx, cfg)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Quotes service should not be nil", func() {
				So(qs, ShouldNotBeNil)
			})

			Convey("When we Started server, it should return no error", func() {
				ctx1, cancel := context.WithCancel(ctx)
				g, gCtx := errgroup.WithContext(ctx1)
				var err1 error

				g.Go(func() error {
					err1 = qs.Startup(gCtx)
					return err1
				})
				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return errors.New("fault")
					case <-time.After(1 * time.Second):
						cancel()
						return nil
					}
				})

				err3 := g.Wait()
				Convey("should be nil", func() {
					So(err3, ShouldBeNil)
				})
			})

			Convey("When starting we can get quotes", func() {
				ctx1, cancel := context.WithCancel(ctx)
				g, gCtx := errgroup.WithContext(ctx1)
				var err1 error

				g.Go(func() error {
					err1 = qs.Startup(gCtx)
					return err1
				})
				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return nil
					case <-time.After(1 * time.Second):
						Convey("When we call for quote", t, func() {
							quote, err3 := qs.GetQuote(gCtx)
							Convey("Should not be error", func() {
								So(err3, ShouldBeNil)
							})
							Convey("Quote should not has zero length", func() {
								So(quote, ShouldNotBeEmpty)
							})
						})
						cancel()
						return nil
					}
				})

				err3 := g.Wait()
				Convey("should not nil", func() {
					So(err3, ShouldBeNil)
				})

			})

		})
	})
}
