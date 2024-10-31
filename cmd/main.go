package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/teyz/tezos-delegation/internal/config"
	database_v1_postgres "github.com/teyz/tezos-delegation/internal/database/v1/postgres"
	handlers_http "github.com/teyz/tezos-delegation/internal/handlers/http"
	poller_v1 "github.com/teyz/tezos-delegation/internal/poller/v1"
	service_v1 "github.com/teyz/tezos-delegation/internal/service/v1"
	pkg_config "github.com/teyz/tezos-delegation/pkg/config"
	"github.com/teyz/tezos-delegation/pkg/database/postgres"
	"github.com/teyz/tezos-delegation/pkg/poller"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	cfg := &config.Config{}
	err := pkg_config.ParseConfig(cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to parse config")
	}

	if err := poller.Validator(&cfg.PollerConfig); err != nil {
		log.Fatal().Err(err).
			Msg("main: invalid poller config")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	databaseConnection, err := postgres.NewDatabaseConnection(ctx, &cfg.PostgresConfig)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create database connection")
	}
	databaseClient := database_v1_postgres.NewClient(ctx, databaseConnection)

	service, err := service_v1.NewService(ctx, databaseClient)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create service")
	}

	poller, err := poller_v1.NewPoller(ctx, service, &cfg.PollerConfig)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create poller")
	}

	if err := poller.StartPolling(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to start poller")
	}

	httpServer, err := handlers_http.NewServer(ctx, cfg.HTTPServerConfig, service)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create http server")
	}

	if err := httpServer.Setup(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to setup http server")
	}

	if err := httpServer.Start(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to start http server")
	}

	<-sigs
	cancel()

	if err := httpServer.Stop(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to stop http server")
	}

	os.Exit(0)
}
