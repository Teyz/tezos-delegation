package database_v1

import (
	"context"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
)

//go:generate mockgen -source interface.go -destination mocks/mock_database.go -package database_mocks
type Database interface {
	//Delegation
	CreateDelegation(ctx context.Context, req *entities_delegation_v1.Delegation_Create) error
	FetchDelegations(ctx context.Context, offset int, limit int) ([]*entities_delegation_v1.Delegation, error)
	FetchDelegationsByYear(ctx context.Context, year int, offset int, limit int) ([]*entities_delegation_v1.Delegation, error)
}
