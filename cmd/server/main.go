package main

import (
	"os"
	"os/signal"
	"syscall"

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

	gzap.Logger.Info("Hello, World")
	printAllEnvVariables()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
