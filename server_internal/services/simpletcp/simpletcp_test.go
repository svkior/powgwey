package simpletcp_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/svkior/powgwey/server_internal/services/simpletcp"
	"golang.org/x/sync/errgroup"
)

const (
	TimeoutOnStartup = 300 * time.Millisecond
)

type configMock struct {
	bindAddr string
}

func (c *configMock) GetBindAddress() string {
	return c.bindAddr
}

type quotesMock struct {
	quote string
	err   error
}

func (q *quotesMock) GetQuote(_ context.Context) (string, error) {
	if q.err != nil {
		return "", q.err
	}
	return q.quote, nil
}

func TestSimpleTCP(t *testing.T) {
	ctx := context.TODO()
	Convey("Given SimpleTCP Server with empty config", t, func() {
		_, err := simpletcp.NewSimpleTCPServer(nil, nil)
		Convey("If no config error should be Nil Config", func() {
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, simpletcp.ErrNilConfig)
		})
	})

	Convey("Given SimpleTCP Server with empty quotes Service", t, func() {
		cfg := &configMock{
			bindAddr: ":0",
		}
		_, err := simpletcp.NewSimpleTCPServer(cfg, nil)
		Convey("If no config error should be Nil Config", func() {
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, simpletcp.ErrNilQuotes)
		})
	})

	Convey("Given SimpleTCP Config", t, func() {
		cfg := &configMock{
			bindAddr: ":0",
		}
		quotes := &quotesMock{
			quote: "Hello, World!",
		}
		Convey("When creating new SimpleTCP server", func() {
			srv, err := simpletcp.NewSimpleTCPServer(cfg, quotes)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
				Convey("The object should not be nil", func() {
					So(srv, ShouldNotBeNil)
					Convey("When starting tcp server", func() {
						ctx1, cancel := context.WithCancel(ctx)
						g, gCtx := errgroup.WithContext(ctx1)
						var err1 error

						g.Go(func() error {
							err1 = srv.Startup(gCtx)
							return err1
						})

						g.Go(func() error {
							select {
							case <-gCtx.Done():
								return errors.New("fail")
							case <-time.After(TimeoutOnStartup):

								Convey("When running tcp server", t, func() {
									port := srv.GetPort()

									Convey("Port should be not 0", func() {
										So(port, ShouldBeGreaterThan, 0)
										connectString := fmt.Sprintf("127.0.0.1:%d", port)
										Convey("When connecting to server", func() {
											conn, err4 := net.Dial("tcp", connectString)
											Convey("Dial should be ok", func() {
												So(err4, ShouldBeNil)

												Convey("Read quote from server", func() {
													reply := make([]byte, 1024)
													bytesRead, err5 := conn.Read(reply)
													Convey("We readed quote from server", func() {
														So(err5, ShouldBeNil)
														So(bytesRead, ShouldEqual, len(quotes.quote))
														gotQuote := string(reply[0:bytesRead])
														So(gotQuote, ShouldEqual, quotes.quote)
													})
												})
												conn.Close()
											})
										})
									})
								})
								cancel()
								return nil
							}
						})

						err2 := g.Wait()
						Convey("startup should not be errornous", func() {
							So(err2, ShouldBeNil)
						})

						Convey("Startup should not be errorful", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
		})
	})
}
