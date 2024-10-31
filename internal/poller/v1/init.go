package poller_v1

import (
	"context"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	service_v1 "github.com/teyz/tezos-delegation/internal/service/v1"
	"github.com/teyz/tezos-delegation/pkg/poller"
)

type Poller struct {
	config  *poller.Config
	service *service_v1.Service
	channel chan *entities_delegation_v1.Delegation_Create
}

func NewPoller(ctx context.Context, service *service_v1.Service, cfg *poller.Config) (*Poller, error) {
	return &Poller{
		config:  cfg,
		service: service,
		channel: make(chan *entities_delegation_v1.Delegation_Create, 25000),
	}, nil
}
