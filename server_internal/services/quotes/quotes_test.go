package quotes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/svkior/powgwey/server_internal/services/quotes"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
)

func TestQuotesSpect(t *testing.T) {
	ctx := context.TODO()
	Convey("Given quotes service with nil config and storage", t, func() {
		_, err := quotes.NewQuotesService(ctx, nil, nil)
		Convey("The error should be ErrNilConfig", func() {
			So(err, ShouldEqual, quotes.ErrNilConfig)
		})
	})

	Convey("Given non empty config for quotes", t, func() {
		cfg := &configMock{
			numberOfWorkers: 1,
		}
		storage := &storageMock{
			quote: "quote",
		}
		Convey("We we create new Quotes Service with zero workers", func() {
			cfg.numberOfWorkers = 0
			_, err := quotes.NewQuotesService(ctx, cfg, storage)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotEqual, nil)
			})
			Convey("Error should be Zero Workers", func() {
				So(err, ShouldEqual, quotes.ErrZeroWorkersCount)
			})
		})

		Convey("When we create new Quotes Service with normal params", func() {
			qs, err := quotes.NewQuotesService(ctx, cfg, storage)

			Convey("Quotes service should not be nil", func() {
				So(err, ShouldBeNil)
				So(qs, ShouldNotBeNil)

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
						case <-time.After(TimeoutOnStartup):
							Convey("When we call for quote", t, func() {
								deadline := time.Now().Add(TimeoutOnStartup)
								runCtx, cancelCtx := context.WithDeadline(gCtx, deadline)
								defer cancelCtx()
								quote, err3 := qs.GetQuote(runCtx)
								Convey("Should not be error", func() {
									So(err3, ShouldBeNil)

									Convey("Quote should not has zero length", func() {
										So(quote, ShouldNotBeEmpty)
									})
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
	})
}
