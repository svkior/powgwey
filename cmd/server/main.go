package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/svkior/powgwey/server_internal/config"

	"gopkg.in/dailymuse/gzap.v1"
)

func printAllEnvVariables() {
	allEnvs := os.Environ()

	for _, value := range allEnvs {
		gzap.Logger.Info("Server",
			gzap.String("env", value),
		)
	}
}

func main() {
	if err := gzap.InitLogger(); err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		gzap.Logger.Fatal("Error loading configuration", gzap.Error(err))
	}

	//FIXME: remove after using
	gzap.Logger.Info("Config from", gzap.String("Service Endpoints HOST", cfg.GetString("HOST")))
	gzap.Logger.Info("Config from", gzap.String("Service Endpoints HOST", cfg.GetString("PORT")))

	gzap.Logger.Info("Hello, World")
	printAllEnvVariables()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
