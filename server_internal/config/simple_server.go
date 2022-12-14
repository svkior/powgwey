package config

import (
	"errors"
	"net"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type simpleServerConfig struct {
	vpr *viper.Viper
}

func (c *simpleServerConfig) GetWorkersCount() uint {
	return c.vpr.GetUint("QUOTES_WORKERS")
}

func (c *simpleServerConfig) GetQuotesFilepath() string {
	return c.vpr.GetString("QUOTES_FILEPATH")
}

func (c *simpleServerConfig) GetProcessingTime() time.Duration {
	return c.vpr.GetDuration("QUOTES_PROCESSING_TIME")
}

func (c *simpleServerConfig) GetBindAddress() string {
	endpoint := net.JoinHostPort(
		c.vpr.GetString("HOST"),
		c.vpr.GetString("PORT"),
	)
	return endpoint
}

func (c *simpleServerConfig) GetReadTimeout() time.Duration {
	return c.vpr.GetDuration("NET_READ_TIMEOUT")
}

func (c *simpleServerConfig) GetWriteTimeout() time.Duration {
	return c.vpr.GetDuration("NET_WRITE_TIMEOUT")
}

func (c *simpleServerConfig) prepare() (err error) {
	c.vpr, err = NewConfig()
	if err != nil {
		return err
	}

	return nil
}

func NewSimpleServerConfig() (*simpleServerConfig, error) {
	cfg := &simpleServerConfig{}

	err := cfg.prepare()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
