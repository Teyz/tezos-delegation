package service_v1

import (
	"context"

	database_v1 "github.com/teyz/tezos-delegation/internal/database/v1"
)

type Service struct {
	store database_v1.Database
}

func NewService(ctx context.Context, store database_v1.Database) (*Service, error) {
	return &Service{
		store: store,
	}, nil
}
