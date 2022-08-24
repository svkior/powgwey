package config_test

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/svkior/powgwey/server_internal/config"
)

func TestSimpleServerConfig(t *testing.T) {
	Convey("Setting some env variables", t, func() {
		os.Setenv("SRV_QUOTES_PROCESSING_TIME", "0.3s")
		os.Setenv("SRV_QUOTES_FILEPATH", "/opt/user/data/movies.json")
		os.Setenv("SRV_QUOTES_WORKERS", "100")
		os.Setenv("SRV_PORT", "8000")
		os.Setenv("SRV_HOST", "0.0.0.0")
		os.Setenv("SRV_NET_READ_TIMEOUT", "1s")
		os.Setenv("SRV_NET_WRITE_TIMEOUT", "1s")
		Convey("When getting new configuration", func() {
			cfg, err := config.NewSimpleServerConfig()
			Convey("There is not errors and  existing config with fields", func() {
				So(err, ShouldBeNil)
				So(cfg, ShouldNotBeNil)
				So(cfg.GetWorkersCount(), ShouldEqual, uint(100))
				So(cfg.GetQuotesFilepath(), ShouldEqual, "/opt/user/data/movies.json")
				So(cfg.GetProcessingTime(), ShouldEqual, 300*time.Millisecond)
				So(cfg.GetBindAddress(), ShouldEqual, "0.0.0.0:8000")
				So(cfg.GetReadTimeout(), ShouldEqual, time.Second)
				So(cfg.GetWriteTimeout(), ShouldEqual, time.Second)
			})
		})
	})
}
