package poller

import "errors"

type Config struct {
	WorkerCount    int    `env:"WORKER_COUNT"`
	PollerInterval int    `env:"POLLER_INTERVAL"`
	PollerUrl      string `env:"POLLER_URL"`
}

func Validator(cfg *Config) error {
	if cfg.WorkerCount <= 0 {
		return errors.New("WORKER_COUNT shouldn't be empty")
	}
	if cfg.PollerInterval <= 0 {
		return errors.New("POLLER_INTERVAL shouldn't be empty")
	}
	if cfg.PollerUrl == "" {
		return errors.New("POLLER_URL shouldn't be empty")
	}
	return nil
}
