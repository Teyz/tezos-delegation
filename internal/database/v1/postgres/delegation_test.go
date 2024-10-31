package database_v1_postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	pkgerrors "github.com/teyz/tezos-delegation/pkg/errors"
)

func Test_CreateDelegation(t *testing.T) {
	t.Run("ok - create delegation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		timestamp := time.Now()

		mock.ExpectExec("INSERT INTO delegations").WithArgs(1, timestamp, 125896, "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", 2338084).WillReturnResult(sqlmock.NewResult(1, 1))

		err = sqlxDB.CreateDelegation(context.Background(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		})
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("nok - create delegation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		timestamp := time.Now()

		mock.ExpectExec("INSERT INTO delegations").WithArgs(1, timestamp, 125896, "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", 2338084).WillReturnError(pkgerrors.NewInternalServerError("error"))

		err = sqlxDB.CreateDelegation(context.Background(), &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: timestamp,
			Amount:    125896,
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     2338084,
		})
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_FetchDelegations(t *testing.T) {
	t.Run("ok - fetch delegations", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		timestamp := time.Now()

		rows := sqlmock.NewRows([]string{"timestamp", "amount", "delegator", "level"}).
			AddRow(timestamp, int64(125896), "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", int64(2338084))

		mock.ExpectQuery("SELECT timestamp, amount, delegator, level FROM delegations ORDER BY timestamp DESC LIMIT $2 OFFSET $1").WillReturnError(nil).WillReturnRows(rows)

		delegations, err := sqlxDB.FetchDelegations(context.Background(), 0, 100)
		assert.NotNil(t, delegations)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(delegations))
		for _, delegation := range delegations {
			assert.Equal(t, timestamp, delegation.Timestamp)
			assert.Equal(t, int64(125896), delegation.Amount)
			assert.Equal(t, "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", delegation.Delegator)
			assert.Equal(t, int64(2338084), delegation.Level)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("nok - fetch delegations", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		mock.ExpectQuery("SELECT timestamp, amount, delegator, level FROM delegations ORDER BY timestamp DESC LIMIT $2 OFFSET $1").WillReturnError(pkgerrors.NewInternalServerError("error"))

		delegations, err := sqlxDB.FetchDelegations(context.Background(), 0, 100)
		assert.Nil(t, delegations)
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func Test_FetchDelegationsByYear(t *testing.T) {
	t.Run("ok - fetch delegations by year", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		timestamp := time.Now()

		rows := sqlmock.NewRows([]string{"timestamp", "amount", "delegator", "level"}).
			AddRow(timestamp, int64(125896), "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", int64(2338084))

		mock.ExpectQuery("SELECT timestamp, amount, delegator, level FROM delegations WHERE EXTRACT(YEAR FROM timestamp) = $1 ORDER BY timestamp DESC LIMIT $3 OFFSET $2").WillReturnError(nil).WillReturnRows(rows)

		delegations, err := sqlxDB.FetchDelegationsByYear(context.Background(), timestamp.Year(), 0, 100)
		assert.NotNil(t, delegations)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(delegations))
		for _, delegation := range delegations {
			assert.Equal(t, timestamp, delegation.Timestamp)
			assert.Equal(t, int64(125896), delegation.Amount)
			assert.Equal(t, "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL", delegation.Delegator)
			assert.Equal(t, int64(2338084), delegation.Level)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("nok - fetch delegations by year", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := &dbClient{
			connection: sqlx.NewDb(db, "sqlmock"),
		}

		mock.ExpectQuery("SELECT timestamp, amount, delegator, level FROM delegations WHERE EXTRACT(YEAR FROM timestamp) = $1 ORDER BY timestamp DESC LIMIT $3 OFFSET $2").WillReturnError(pkgerrors.NewInternalServerError("error"))

		delegations, err := sqlxDB.FetchDelegationsByYear(context.Background(), time.Now().Year(), 0, 100)
		assert.Nil(t, delegations)
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
