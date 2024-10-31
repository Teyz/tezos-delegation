package http

import (
	"context"
	"net/http"

	pkgerrors "github.com/teyz/tezos-delegation/pkg/errors"
)

func TranslateError(ctx context.Context, err error) int {
	switch {
	case pkgerrors.IsNotFoundError(err):
		return http.StatusNotFound
	case pkgerrors.IsResourceAlreadyCreatedError(err):
		return http.StatusConflict
	case pkgerrors.IsBadRequestError(err):
		return http.StatusBadRequest
	case pkgerrors.IsUnauthorizedError(err):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
