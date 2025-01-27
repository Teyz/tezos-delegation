package handlers_http_delegations_v1

import (
	"context"

	service_v1 "github.com/teyz/tezos-delegation/internal/service/v1"
)

type Handler struct {
	service *service_v1.Service
}

func NewHandler(_ context.Context, service *service_v1.Service) *Handler {
	return &Handler{
		service: service,
	}
}
