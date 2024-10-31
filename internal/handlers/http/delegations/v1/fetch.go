package handlers_http_delegations_v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	entities_delegation_v1 "github.com/teyz/tezos-delegation/internal/entities/delegation/v1"
	pkg_http "github.com/teyz/tezos-delegation/pkg/http"
)

type FetchDelegationsResponse struct {
	Data []*entities_delegation_v1.Delegation `json:"data"`
}

func (h *Handler) FetchDelegations(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	year := req.URL.Query().Get("year")
	offset := req.URL.Query().Get("offset")
	limit := req.URL.Query().Get("limit")

	if offset == "" {
		offset = "0"
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		log.Error().Err(err).
			Msg("handlers.http.delegations.v1.fetch.Handler.FetchDelegations: can not parse offset")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if limit == "" {
		limit = "100"
	}

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Error().Err(err).
			Msg("handlers.http.delegations.v1.fetch.Handler.FetchDelegations: can not parse limit")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var delegations []*entities_delegation_v1.Delegation

	if year != "" {
		parsedYear, err := strconv.Atoi(year)
		if err != nil {
			log.Error().Err(err).
				Msg("handlers.http.delegations.v1.fetch.Handler.FetchDelegations: can not parse year")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		delegations, err = h.service.FetchDelegationsByYear(ctx, parsedYear, parsedOffset, parsedLimit)
		if err != nil {
			http.Error(w, err.Error(), pkg_http.TranslateError(ctx, err))
		}
	} else {
		delegations, err = h.service.FetchDelegations(ctx, parsedOffset, parsedLimit)
		if err != nil {
			http.Error(w, err.Error(), pkg_http.TranslateError(ctx, err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&FetchDelegationsResponse{Data: delegations})
}
