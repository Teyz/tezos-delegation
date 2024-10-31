package config

import "github.com/caarlos0/env/v10"

type Config struct {
	ServiceName string `env:"SERVICE_NAME"`
	Environment string `env:"ENVIRONMENT" envDefault:"local"`
}

func ParseConfig[T any](cfg *T) error {
	return env.Parse(cfg)
}
