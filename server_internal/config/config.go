package config

import (
	"github.com/svkior/powgwey/server_internal/app"

	"github.com/spf13/viper"
	"gopkg.in/dailymuse/gzap.v1"
)

func NewConfig() (*viper.Viper, error) {
	// Setup Config
	cfg := viper.New()

	// Load Config
	cfg.SetEnvPrefix(app.ConfigPrefix)
	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	err := cfg.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			gzap.Logger.Warn("Error when Fetching Configuration",
				gzap.Error(err),
			)
			return nil, err
		}
	}

	return cfg, nil
}
