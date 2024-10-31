package service_v1

import (
	"context"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
)

type KilnService interface {
	CreateDelegation(ctx context.Context, delegation *entities_delegation_v1.Delegation_Create) error
	FetchDelegations(ctx context.Context, offset int, limit int) ([]*entities_delegation_v1.Delegation, error)
	FetchDelegationsByYear(ctx context.Context, year int, offset int, limit int) ([]*entities_delegation_v1.Delegation, error)
}
