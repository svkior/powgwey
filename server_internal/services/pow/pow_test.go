package pow_test

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/svkior/powgwey/server_internal/models"
	"github.com/svkior/powgwey/server_internal/services/pow"
	"golang.org/x/sync/errgroup"
)

const clientTimeout = time.Millisecond * 300

var (
	ErrTestTimeout = errors.New("test timeout")
	ErrNoError     = errors.New("no error")
)

type cfgMock struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (m *cfgMock) GetReadTimeout() time.Duration {
	return m.readTimeout
}

func (m *cfgMock) GetWriteTimeout() time.Duration {
	return m.writeTimeout
}

type mertricMock struct {
	diff uint
}

func (m *mertricMock) GetDifficulties() uint {
	return m.diff
}

func TestPoWMiddleware(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	ctx := context.TODO()

	Convey("When get PoW middleware with empty cfg", t, func() {
		_, err := pow.NewPoWMiddleware(nil, nil)
		Convey("Should  be error", func() {
			So(err, ShouldResemble, pow.ErrNilConfig)
		})
	})

	cfg := &cfgMock{
		readTimeout:  100 * time.Millisecond,
		writeTimeout: 100 * time.Millisecond,
	}

	Convey("When get PoW middleware with empte metrics", t, func() {
		_, err := pow.NewPoWMiddleware(cfg, nil)
		Convey("Should  be error", func() {
			So(err, ShouldResemble, pow.ErrNilMetrics)
		})
	})

	Convey("When get PoW middleware", t, func() {
		mm := &mertricMock{
			diff: 10,
		}
		m, err := pow.NewPoWMiddleware(cfg, mm)
		Convey("Should be no error", func() {
			So(err, ShouldBeNil)
			So(m, ShouldNotBeNil)

			Convey("We can't Validate with empty connection", func() {
				err1 := m.Validate(ctx, nil)

				Convey("Error should be nil connection", func() {
					So(err1, ShouldResemble, pow.ErrNilConnection)
				})
			})

			Convey("We can't Validate conneciton with timeout", func() {
				server, _ := net.Pipe()
				lCtx, lCancel := context.WithDeadline(ctx, time.Now().Add(clientTimeout))
				var err3 error
				var noTimeout bool
				select {
				case e := <-func() chan error {
					c := make(chan error, 1)
					e := m.Validate(lCtx, server)
					c <- e
					return c
				}():
					err3 = e
				case <-lCtx.Done():
					noTimeout = true
				}
				lCancel()
				Convey("Validate should return timeout error", func() {
					So(err3, ShouldResemble, pow.ErrTimeout)
					So(noTimeout, ShouldBeFalse)
				})
			})

			Convey("When we validate with wrong request", func() {
				server, client := net.Pipe()
				lCtx, lCancel := context.WithDeadline(ctx, time.Now().Add(clientTimeout))
				var err3 error

				g, gCtx := errgroup.WithContext(lCtx)

				g.Go(func() error {
					select {
					case e := <-func() chan error {
						c := make(chan error, 1)
						e := m.Validate(gCtx, server)
						c <- e
						return c
					}():
						if e == nil {
							return ErrNoError
						}

						return e
					case <-lCtx.Done():
						return ErrTestTimeout
					}

				})
				_, err3 = client.Write([]byte("Simple answer"))
				Convey("err  for write should be nil", func() {
					So(err3, ShouldBeNil)
				})
				err3 = g.Wait()
				lCancel()
				Convey("Error after wait should be nil", func() {
					So(err3, ShouldResemble, pow.ErrGetRequest)
				})
			})
		})
	})
}

func TestMiddlewareComm1(t *testing.T) {
	ctx := context.TODO()
	Convey("When get PoW middleware", t, func() {
		mm := &mertricMock{
			diff: 10,
		}
		cfg := &cfgMock{
			readTimeout:  100 * time.Millisecond,
			writeTimeout: 100 * time.Millisecond,
		}

		m, err := pow.NewPoWMiddleware(cfg, mm)
		Convey("Should be no error", func() {
			So(err, ShouldBeNil)
			So(m, ShouldNotBeNil)
			Convey("When we validate with ok request, timeout on solve", func() {
				server, client := net.Pipe()
				lCtx, lCancel := context.WithDeadline(ctx, time.Now().Add(clientTimeout))
				var err3 error
				g, gCtx := errgroup.WithContext(lCtx)

				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return nil
					case <-time.After(300 * time.Millisecond):
						lCancel()
						return ErrTestTimeout
					}
				})

				g.Go(func() error {
					return m.Validate(gCtx, server)
				})

				g.Go(func() error {
					reply := make([]byte, 1024)
					req := models.Request{}
					bytes_marshalled := req.MarshalTo(reply)
					_, err3 = client.Write(reply[0:bytes_marshalled])
					Convey("Write error  should be nil", t, func() {
						So(err3, ShouldBeNil)
						bytes_readed, err3 := client.Read(reply)
						Convey("Read error should be nil", func() {

							So(err3, ShouldBeNil)
							Convey("Read error  should be nil", func() {
								So(err3, ShouldBeNil)
								puzzle := &models.Puzzle{}
								_, err3 = puzzle.Unmarshal(reply[0:bytes_readed])
								Convey("Complete unmarshal puzzle", func() {
									So(err3, ShouldBeNil)
									solution, err3 := pow.Solution(puzzle, pow.Algosin)
									Convey("We should have solution", func() {
										So(err3, ShouldBeNil)
										bytes_marshalled := solution.MarshalTo(reply)
										_, err3 = client.Write(reply[0:bytes_marshalled])
										Convey("Write error  should be nil", func() {
											So(err3, ShouldBeNil)
											bytes_readed, err3 = client.Read(reply)
											Convey("Read error should be nil", func() {
												So(err3, ShouldBeNil)
												ok := &models.Result{}
												_, err3 = ok.Unmarshal(reply[0:bytes_readed])
												Convey("Unmarshal and OK", func() {
													So(err3, ShouldBeNil)
													So(ok.Ok, ShouldBeTrue)
												})
											})
										})
									})

								})
							})
						})
					})
					return ErrNoError
				})

				err3 = g.Wait()
				Convey("Error after wait should be nil", func() {
					So(err3, ShouldResemble, ErrNoError)
				})
			})

		})
	})
}
