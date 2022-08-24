package main

import (
	"os"

	"gopkg.in/dailymuse/gzap.v1"
)

func printAllEnvVariables() {
	allEnvs := os.Environ()

	for _, value := range allEnvs {
		gzap.Logger.Info("Client",
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

}
