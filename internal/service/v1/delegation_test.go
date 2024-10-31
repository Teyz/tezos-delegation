package service_v1

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	database_mocks "github.com/teyz/tezos-delegation/internal/database/v1/mocks"
	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	pkgerrors "github.com/teyz/tezos-delegation/pkg/errors"
)

func Test_CreateDelegation(t *testing.T) {
	t.Run("ok - create delegation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		timestamp := time.Now()

		mock_database.EXPECT().CreateDelegation(gomock.Any(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		}).Return(nil)

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		err = s.CreateDelegation(context.Background(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		})
		assert.NoError(t, err)
	})
	t.Run("nok - create delegation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		timestamp := time.Now()

		mock_database.EXPECT().CreateDelegation(gomock.Any(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		}).Return(pkgerrors.NewInternalServerError("error"))

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		err = s.CreateDelegation(context.Background(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		})
		assert.Error(t, err)
	})
}

func Test_FetchDelegations(t *testing.T) {
	t.Run("ok - fetch delegations from database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		today := time.Now()
		lastYear := today.AddDate(-1, 0, 0)

		delegationsExpected := []*entities_delegation_v1.Delegation{
			{
				Timestamp: today,
				Amount:    125896,
				Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
				Level:     2338084,
			},
			{
				Timestamp: lastYear,
				Amount:    9856354,
				Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
				Level:     1461334,
			},
		}

		mock_database.EXPECT().FetchDelegations(gomock.Any(), 0, 100).Return(delegationsExpected, nil)

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		delegations, err := s.FetchDelegations(context.Background(), 0, 100)
		assert.NotNil(t, delegations)
		assert.NoError(t, err)

		for _, delegation := range delegations {
			assert.Contains(t, delegationsExpected, delegation)
		}
	})
	t.Run("nok - fetch delegations from database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		mock_database.EXPECT().FetchDelegations(gomock.Any(), 0, 100).Return(nil, pkgerrors.NewNotFoundError("error"))

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		delegations, err := s.FetchDelegations(context.Background(), 0, 100)
		assert.Nil(t, delegations)
		assert.Error(t, err)
	})
}

func Test_FetchDelegationsByYear(t *testing.T) {
	t.Run("ok - fetch delegations by year from database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		today := time.Now()
		lastYear := today.AddDate(-1, 0, 0)

		delegationsAdded := []*entities_delegation_v1.Delegation_Create{
			{
				ID:        2,
				Timestamp: today,
				Amount:    125896,
				Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
				Level:     2338084,
			},
			{
				ID:        1,
				Timestamp: lastYear,
				Amount:    9856354,
				Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
				Level:     1461334,
			},
		}

		delegationsExpected := []*entities_delegation_v1.Delegation{
			{
				Timestamp: today,
				Amount:    125896,
				Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
				Level:     2338084,
			},
		}

		mock_database.EXPECT().FetchDelegationsByYear(gomock.Any(), today.Year(), 0, 100).Return(delegationsExpected, nil)

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		for _, delegation := range delegationsAdded {
			mock_database.EXPECT().CreateDelegation(gomock.Any(), delegation).Return(nil)

			err := s.CreateDelegation(context.Background(), delegation)
			assert.NoError(t, err)
		}

		delegations, err := s.FetchDelegationsByYear(context.Background(), today.Year(), 0, 100)
		assert.NotNil(t, delegations)
		assert.NoError(t, err)

		assert.Len(t, delegations, 1)
		for _, delegation := range delegations {
			assert.Contains(t, delegationsExpected, delegation)
		}
	})
	t.Run("nok - fetch delegations by year from database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		today := time.Now()

		mock_database.EXPECT().FetchDelegationsByYear(gomock.Any(), today.Year(), 0, 100).Return(nil, pkgerrors.NewNotFoundError("error"))

		s, err := NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		delegations, err := s.FetchDelegationsByYear(context.Background(), today.Year(), 0, 100)
		assert.Nil(t, delegations)
		assert.Error(t, err)
	})
}
