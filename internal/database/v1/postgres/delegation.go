package database_v1_postgres

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	"github.com/teyz/tezos-delegation/pkg/errors"
)

func (d *dbClient) CreateDelegation(ctx context.Context, req *entities_delegation_v1.Delegation_Create) error {
	_, err := d.connection.DB.ExecContext(ctx,
		`INSERT INTO 
			delegations (
				id,
				timestamp,
				amount,
				delegator,
				level
			) 
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING;
		`,
		req.ID, req.Timestamp, req.Amount, req.Delegator, req.Level)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.pgx.dbClient.CreateDelegation: failed to create delegation: %v", err.Error())
		return errors.NewInternalServerError(fmt.Sprintf("database.pgx.dbClient.CreateDelegation: failed to create delegation: %v", err.Error()))
	}

	return nil
}

func (d *dbClient) FetchDelegations(ctx context.Context, offset int, limit int) ([]*entities_delegation_v1.Delegation, error) {
	rows, err := d.connection.DB.QueryContext(ctx, `
		SELECT
			timestamp,
			amount,
			delegator,
			level
		FROM
			delegations
		ORDER BY timestamp DESC
		LIMIT $2
		OFFSET $1
	`, offset, limit)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.postgres.dbClient.FetchDelegations: failed to fetch delegations: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchDelegations: failed to fetch delegations: %v", err.Error()))
	}
	defer rows.Close()

	delegations := make([]*entities_delegation_v1.Delegation, 0)

	for rows.Next() {
		delegation := &entities_delegation_v1.Delegation{}

		err := rows.Scan(
			&delegation.Timestamp,
			&delegation.Amount,
			&delegation.Delegator,
			&delegation.Level,
		)
		if err != nil {
			log.Error().Err(err).
				Msgf("database.postgres.dbClient.FetchDelegations: failed to scan delegation: %v", err.Error())
			return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchDelegations: failed to scan delegation: %v", err.Error()))
		}

		delegations = append(delegations, delegation)
	}

	return delegations, nil
}

func (d *dbClient) FetchDelegationsByYear(ctx context.Context, year int, offset int, limit int) ([]*entities_delegation_v1.Delegation, error) {
	rows, err := d.connection.DB.QueryContext(ctx, `
		SELECT
			timestamp,
			amount,
			delegator,
			level
		FROM
			delegations
		WHERE
			EXTRACT(YEAR FROM timestamp) = $1
		ORDER BY timestamp DESC
		LIMIT $3
		OFFSET $2
	`, year, offset, limit)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.postgres.dbClient.FetchDelegationsByYear: failed to fetch delegations: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchDelegationsByYear: failed to fetch delegations: %v", err.Error()))
	}
	defer rows.Close()

	delegations := make([]*entities_delegation_v1.Delegation, 0)

	for rows.Next() {
		delegation := &entities_delegation_v1.Delegation{}

		err := rows.Scan(
			&delegation.Timestamp,
			&delegation.Amount,
			&delegation.Delegator,
			&delegation.Level,
		)
		if err != nil {
			log.Error().Err(err).
				Msgf("database.postgres.dbClient.FetchDelegationsByYear: failed to scan delegation: %v", err.Error())
			return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchDelegationsByYear: failed to scan delegation: %v", err.Error()))
		}

		delegations = append(delegations, delegation)
	}

	return delegations, nil
}
