package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/app"
	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/server/http"
	// sqlstorage "github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage/sql".
	memorystorage "github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(configFile)
	// config := NewConfig("configs/config.toml")
	logg := logger.New(config.Logger.Level)

	// storage := sqlstorage.New(ctx, config.Storage.postgresDSN)
	storage := memorystorage.New(ctx)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx, config.Server.Host, config.Server.HTTPPort); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		return
	}
}
