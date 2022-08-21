package quotes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/svkior/powgwey/server_internal/models"
	"github.com/svkior/powgwey/server_internal/services/quotes"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
)

func TestWorkerSpect(t *testing.T) {
	ctx := context.TODO()
	Convey("Given Worker with nil storage", t, func() {
		_, err := quotes.NewQuotesWorker(nil, nil)
		Convey("Error should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
		Convey("Error should be nil storage", func() {
			So(err, ShouldResemble, quotes.ErrStorageIsNil)
		})
	})

	Convey("Given Worker with nil workerChannel", t, func() {
		storage := &storageMock{
			quote: "quote",
		}
		_, err := quotes.NewQuotesWorker(storage, nil)
		Convey("Error should be nil storage", func() {
			So(err, ShouldResemble, quotes.ErrNilWorkerChannel)
		})
	})

	Convey("Given mock storage and workerChannel", t, func() {
		storage := &storageMock{
			quote: "quote",
		}
		workerChannel := make(chan chan *models.QuotesWork, 10)

		Convey("When create worker", func() {
			w, err := quotes.NewQuotesWorker(storage, workerChannel)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)

				Convey("Worker should not be nil", func() {
					So(w, ShouldNotBeNil)
				})

				Convey("When starting worker", func() {
					ctx1, cancel := context.WithCancel(ctx)
					g, gCtx := errgroup.WithContext(ctx1)
					var err1 error

					g.Go(func() error {
						err1 = w.Startup(gCtx)
						return err1
					})

					g.Go(func() error {
						select {
						case <-gCtx.Done():
							return errors.New("fault")
						case <-time.After(TimeoutOnStartup):

							Convey("When we have a job", t, func() {
								finishChannel := make(chan struct{}, 1)
								job := models.QuotesWork{
									Finish: finishChannel,
								}

								Convey("When we Call the job directly", func() {
									err4 := w.DoWork(gCtx, &job)
									Convey("Work Should be called without error", func() {
										So(err4, ShouldBeNil)
										Convey("When we wait for job answer", func() {
											select {
											case <-time.After(1 * time.Second):
												Convey("We Have TIMEOUT", func() {
													So(nil, ShouldNotBeNil)
												})
											case <-finishChannel:
												Convey("Return from worker should not contains error", func() {
													So(job.Error, ShouldBeNil)
													So(job.Quote, ShouldEqual, "quote")
												})
											}
										})
									})
								})
								Convey("When we Call the job through channel", func() {
									var wrk chan *models.QuotesWork
									select {
									case <-time.After(TimeoutOnStartup):
										panic("timeout")
									case wrk = <-workerChannel:
										Convey("Work should not be nil", func() {
											So(wrk, ShouldNotBeNil)
											Convey("Send job to the worker", func() {
												job.Quote = ""
												storage.quote = "quote1"
												select {
												case wrk <- &job:
												default:
													panic("cant put job to the channel")
												}
												select {
												case <-time.After(TimeoutOnStartup):
													Convey("We Have TIMEOUT", func() {
														panic("timeout")
													})
												case <-finishChannel:
													Convey("Return from worker should not contains error", func() {
														So(job.Error, ShouldBeNil)
														So(job.Quote, ShouldEqual, "quote1")
													})
												}
											})
										})
									}
								})
							})
							cancel()
							return nil
						}
					})
					err3 := g.Wait()
					Convey("The worker should start without errors", func() {
						So(err3, ShouldBeNil)
					})
				})
			})
		})
	})
}
