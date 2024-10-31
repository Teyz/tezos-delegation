package poller_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
)

func (p *Poller) StartPolling(ctx context.Context) error {
	log.Info().
		Msg("poller.v1.Poller.StartPolling: start polling for delegations")

	go p.PollingHistoricalDelegations(ctx, 0)
	go p.PollingNewDelegations(ctx)

	for range p.config.WorkerCount {
		go p.ConsumeDelegation(ctx)
	}

	return nil
}

func (p *Poller) PollingNewDelegations(ctx context.Context) {
	var wg sync.WaitGroup

	ticker := time.NewTicker(time.Duration(p.config.PollerInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().
				Msg("poller.v1.Poller.PollingNewDelegations: finishing polling for delegations")
			return
		case <-ticker.C:
			wg.Add(1)

			log.Info().
				Msg("poller.v1.Poller.PollingNewDelegations: polling new delegations")

			go func() {
				defer wg.Done()

				resp, err := http.Get(fmt.Sprintf("%v&sort.desc=level", p.config.PollerUrl))
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.PollingNewDelegations: failed to get: %v", err.Error())
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.PollingNewDelegations: failed to read: %v", err.Error())
					return
				}

				var delegations []*entities_delegation_v1.TZKTDelegation
				err = json.Unmarshal(body, &delegations)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.PollingNewDelegations: failed to unmarshal: %v", err.Error())
					return
				}

				for _, d := range delegations {
					p.channel <- &entities_delegation_v1.Delegation_Create{
						ID:        d.ID,
						Timestamp: d.Timestamp,
						Amount:    d.Amount,
						Delegator: d.Sender.Address,
						Level:     d.Level,
					}
				}
			}()
		}
		wg.Wait()
	}
}

func (p *Poller) ConsumeDelegation(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			log.Info().
				Msg("poller.v1.Poller.ConsumeDelegation: finishing polling for delegations")
			return
		case delegation := <-p.channel:
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := p.service.CreateDelegation(ctx, delegation)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.ConsumeDelegation: failed to create delegation: %v", err.Error())
					return
				}
			}()
		}
		wg.Wait()
	}
}

func (p *Poller) PollingHistoricalDelegations(ctx context.Context, cursor int64) {
	var wg sync.WaitGroup

	isEnd := false

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Info().
		Msg("poller.v1.Poller.getHistoricalDelegations: starting polling historical delegations")

	for !isEnd {
		select {
		case <-ctx.Done():
			log.Info().
				Msg("poller.v1.Poller.getHistoricalDelegations: finishing polling for delegations")
			return
		case <-ticker.C:
			wg.Add(1)

			go func(currentCursor *int64, isEnd *bool) {
				defer wg.Done()
				client := http.Client{
					Timeout: 30 * time.Second,
				}
				req, err := http.NewRequest("GET", p.config.PollerUrl, nil)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.PollingHistoricalDelegations: failed to create request: %v", err.Error())
					return
				}
				q := req.URL.Query()
				q.Add("limit", "10000")
				if *currentCursor != 0 {
					q.Add("offset.cr", strconv.FormatInt(*currentCursor, 10))
				}
				req.URL.RawQuery = q.Encode()

				resp, err := client.Do(req)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.PollingHistoricalDelegations: failed to get: %v", err.Error())
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.getHistoricalDelegations: failed to read: %v", err.Error())
					return
				}

				var delegations []*entities_delegation_v1.TZKTDelegation
				err = json.Unmarshal(body, &delegations)
				if err != nil {
					log.Error().Err(err).
						Msgf("poller.v1.Poller.getHistoricalDelegations: failed to unmarshal: %v", err.Error())
					return
				}

				for _, d := range delegations {
					p.channel <- &entities_delegation_v1.Delegation_Create{
						ID:        d.ID,
						Timestamp: d.Timestamp,
						Amount:    d.Amount,
						Delegator: d.Sender.Address,
						Level:     d.Level,
					}
				}

				if len(delegations) < 10000 {
					*isEnd = true
					return
				}

				*currentCursor = delegations[len(delegations)-1].ID

			}(&cursor, &isEnd)
		}
	}

	log.Info().
		Msg("poller.v1.Poller.getHistoricalDelegations: finishing polling historical delegations")

	wg.Wait()
}
