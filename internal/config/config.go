package config

import (
	"github.com/teyz/tezos-delegation/pkg/config"
	"github.com/teyz/tezos-delegation/pkg/database/postgres"
	"github.com/teyz/tezos-delegation/pkg/http"
	"github.com/teyz/tezos-delegation/pkg/poller"
)

type Config struct {
	ServiceConfig config.Config

	HTTPServerConfig http.HTTPServerConfig
	PostgresConfig   postgres.Config
	PollerConfig     poller.Config
}
