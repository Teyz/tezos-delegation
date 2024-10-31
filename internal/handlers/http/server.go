package handlers_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/teyz/tezos-delegation/internal/handlers"
	handlers_http_delegations_v1 "github.com/teyz/tezos-delegation/internal/handlers/http/delegations/v1"
	service_v1 "github.com/teyz/tezos-delegation/internal/service/v1"
	pkp_http "github.com/teyz/tezos-delegation/pkg/http"
)

type httpServer struct {
	router  *http.ServeMux
	config  pkp_http.HTTPServerConfig
	service *service_v1.Service
}

func NewServer(ctx context.Context, cfg pkp_http.HTTPServerConfig, service *service_v1.Service) (handlers.Server, error) {
	return &httpServer{
		router:  http.NewServeMux(),
		config:  cfg,
		service: service,
	}, nil
}

func (s *httpServer) Setup(ctx context.Context) error {
	log.Info().
		Msg("handlers.http.httpServer.Setup: Setting up HTTP server...")

	internalDelegationsV1Handlers := handlers_http_delegations_v1.NewHandler(ctx, s.service)

	s.router.HandleFunc("GET /xtz/delegations", internalDelegationsV1Handlers.FetchDelegations)

	return nil
}

func (s *httpServer) Start(ctx context.Context) error {
	log.Info().
		Uint16("port", s.config.Port).
		Msg("handlers.http.httpServer.Start: Starting HTTP server...")

	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}

func (s *httpServer) Stop(ctx context.Context) error {
	log.Info().
		Msg("handlers.http.httpServer.Stop: Stopping HTTP server...")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.router,
	}

	return server.Shutdown(ctx)
}
