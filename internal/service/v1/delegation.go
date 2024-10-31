package service_v1

import (
	"context"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
)

func (s *Service) CreateDelegation(ctx context.Context, delegation *entities_delegation_v1.Delegation_Create) error {
	err := s.store.CreateDelegation(ctx, delegation)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FetchDelegations(ctx context.Context, offset int, limit int) ([]*entities_delegation_v1.Delegation, error) {
	delegations, err := s.store.FetchDelegations(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return delegations, nil
}

func (s *Service) FetchDelegationsByYear(ctx context.Context, year int, offset int, limit int) ([]*entities_delegation_v1.Delegation, error) {
	delegations, err := s.store.FetchDelegationsByYear(ctx, year, offset, limit)
	if err != nil {
		return nil, err
	}

	return delegations, nil
}
