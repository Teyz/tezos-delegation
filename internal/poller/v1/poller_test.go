package poller_v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	database_mocks "github.com/teyz/tezos-delegation/internal/database/v1/mocks"
	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	service_v1 "github.com/teyz/tezos-delegation/internal/service/v1"
	pkgerrors "github.com/teyz/tezos-delegation/pkg/errors"
	"github.com/teyz/tezos-delegation/pkg/poller"
)

func Test_PollingNewDelegations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		delegations := []*entities_delegation_v1.TZKTDelegation{
			{ID: 1, Timestamp: time.Now(), Amount: 100, Sender: &entities_delegation_v1.TZKTSender{Address: "tz1"}, Level: 100},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(delegations)
	}))
	defer server.Close()

	p := &Poller{
		config: &poller.Config{
			PollerUrl:      server.URL + "?sort.desc=level",
			PollerInterval: 1,
			WorkerCount:    1,
		},
		channel: make(chan *entities_delegation_v1.Delegation_Create, 10),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p.PollingNewDelegations(ctx)

	time.Sleep(2 * time.Second)

	select {
	case delegation := <-p.channel:
		assert.Equal(t, "tz1", delegation.Delegator)
		assert.Equal(t, int64(100), delegation.Amount)
	default:
		t.Fatal("Aucune délégation n'a été ajoutée au channel")
	}
}

func Test_ConsumeDelegation(t *testing.T) {
	t.Run("ok - consume delegation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		delegation := &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: time.Now(),
			Amount:    100,
			Delegator: "tz1",
			Level:     100,
		}

		mock_database.EXPECT().CreateDelegation(gomock.Any(), delegation).Return(nil)

		s, err := service_v1.NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		p := &Poller{
			channel: make(chan *entities_delegation_v1.Delegation_Create, 1),
			service: s,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p.channel <- delegation

		go p.ConsumeDelegation(ctx)

		time.Sleep(1 * time.Second)
	})
	t.Run("nok - consume delegation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock_database := database_mocks.NewMockDatabase(ctrl)

		delegation := &entities_delegation_v1.Delegation_Create{
			ID:        1,
			Timestamp: time.Now(),
			Amount:    100,
			Delegator: "tz1",
			Level:     100,
		}

		mock_database.EXPECT().CreateDelegation(gomock.Any(), delegation).Return(pkgerrors.NewInternalServerError("error"))

		s, err := service_v1.NewService(context.Background(), mock_database)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		p := &Poller{
			channel: make(chan *entities_delegation_v1.Delegation_Create, 1),
			service: s,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p.channel <- delegation

		go p.ConsumeDelegation(ctx)

		time.Sleep(1 * time.Second)
	})
}

func Test_PollingHistoricalDelegations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		delegations := []*entities_delegation_v1.TZKTDelegation{
			{ID: 1, Timestamp: time.Now(), Amount: 100, Sender: &entities_delegation_v1.TZKTSender{Address: "tz1"}, Level: 100},
			{ID: 2, Timestamp: time.Now(), Amount: 200, Sender: &entities_delegation_v1.TZKTSender{Address: "tz2"}, Level: 200},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(delegations)
	}))
	defer server.Close()

	p := &Poller{
		config: &poller.Config{
			PollerUrl:      server.URL + "?limit=10000",
			PollerInterval: 1,
			WorkerCount:    1,
		},
		channel: make(chan *entities_delegation_v1.Delegation_Create, 10),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		p.PollingHistoricalDelegations(ctx, 0)
		wg.Done()
	}()

	wg.Wait()

	time.Sleep(1 * time.Second)

	select {
	case delegation := <-p.channel:
		assert.Equal(t, "tz1", delegation.Delegator)
		assert.Equal(t, int64(100), delegation.Amount)
	default:
		t.Fatal("Aucune délégation n'a été ajoutée au channel")
	}
}
