package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/svkior/powgwey/server_internal/config"
	"github.com/svkior/powgwey/server_internal/services/metrics"
	"github.com/svkior/powgwey/server_internal/services/pow"
	"github.com/svkior/powgwey/server_internal/services/quotes"
	"github.com/svkior/powgwey/server_internal/services/simpletcp"
	"github.com/svkior/powgwey/server_internal/storage"
	"golang.org/x/sync/errgroup"

	"gopkg.in/dailymuse/gzap.v1"
)

func PrintAllEnvVariables() {
	allEnvs := os.Environ()

	for _, value := range allEnvs {
		gzap.Logger.Info("Server",
			gzap.String("env", value),
		)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	if err := gzap.InitLogger(); err != nil {
		panic(err)
	}
	//PrintAllEnvVariables()

	cfg, err := config.NewSimpleServerConfig()
	if err != nil {
		gzap.Logger.Fatal("Error loading configuration", gzap.Error(err))
	}

	// Storage
	store, err := storage.NewQuotesStorage(cfg)
	if err != nil {
		gzap.Logger.Fatal("Error ctreating storage", gzap.Error(err))
	}
	// Quotes
	quotes, err := quotes.NewQuotesService(cfg, store)
	if err != nil {
		gzap.Logger.Fatal("Error ctreating quotes service", gzap.Error(err))
	}

	// Metrics
	metric, err := metrics.NewMetricService(cfg)
	if err != nil {
		gzap.Logger.Fatal("Error ctreating metric service", gzap.Error(err))
	}

	//  PoW middleware
	powMiddleware, err := pow.NewPoWMiddleware(cfg, metric)
	if err != nil {
		gzap.Logger.Fatal("Error ctreating PoW middleWare", gzap.Error(err))
	}

	// TCP Server
	tcpServer, err := simpletcp.NewSimpleTCPServer(cfg, quotes, powMiddleware)
	if err != nil {
		gzap.Logger.Fatal("Error ctreating simple TCP server", gzap.Error(err))
	}

	g, gCtx := errgroup.WithContext(ctx)

	// Storage
	g.Go(func() error {
		return store.Startup(gCtx)
	})
	// Quotes
	g.Go(func() error {
		return quotes.Startup(gCtx)
	})
	// Metrics
	//  Not Need

	//  PoW middleware
	//  Not Need

	// TCP Server
	g.Go(func() error {
		return tcpServer.Startup(gCtx)
	})

	// Starting the waiting goroutine for terminate goroutine
	g.Go(func() error {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		gzap.Logger.Warn("Stoping service because of SIGTERM")
		cancel()
		return nil
	})

	err = g.Wait()
	if err != nil {
		gzap.Logger.Error("Stoping service because errors", gzap.Error(err))
	}

}
