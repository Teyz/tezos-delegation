package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pkgerrors "github.com/teyz/tezos-delegation/pkg/errors"
)

func Test_Errors(t *testing.T) {
	notFoundError := pkgerrors.NewNotFoundError("not found")
	assert.True(t, pkgerrors.IsNotFoundError(notFoundError))
	assert.False(t, pkgerrors.IsBadRequestError(notFoundError))
	assert.Equal(t, "not found", notFoundError.Error())

	badRequestError := pkgerrors.NewBadRequestError("bad request")
	assert.True(t, pkgerrors.IsBadRequestError(badRequestError))
	assert.False(t, pkgerrors.IsInternalServerError(badRequestError))
	assert.Equal(t, "bad request", badRequestError.Error())

	expiredResourceError := pkgerrors.NewExpiredResourceError("expired resource")
	assert.True(t, pkgerrors.IsExpiredResourceError(expiredResourceError))
	assert.False(t, pkgerrors.IsNotFoundError(expiredResourceError))
	assert.Equal(t, "expired resource", expiredResourceError.Error())

	internalServerError := pkgerrors.NewInternalServerError("internal error")
	assert.True(t, pkgerrors.IsInternalServerError(internalServerError))
	assert.False(t, pkgerrors.IsBadRequestError(internalServerError))
	assert.Equal(t, "internal error", internalServerError.Error())

	unauthorizedError := pkgerrors.NewUnauthorizedError("unauthorized")
	assert.True(t, pkgerrors.IsUnauthorizedError(unauthorizedError))
	assert.False(t, pkgerrors.IsInternalServerError(unauthorizedError))
	assert.Equal(t, "unauthorized", unauthorizedError.Error())

	alreadyCreatedError := pkgerrors.NewResourceAlreadyCreatedError("already created")
	assert.True(t, pkgerrors.IsResourceAlreadyCreatedError(alreadyCreatedError))
	assert.False(t, pkgerrors.IsNotFoundError(alreadyCreatedError))
	assert.Equal(t, "already created", alreadyCreatedError.Error())

	outdatedError := pkgerrors.NewOutdatedResourceError("outdated")
	assert.True(t, pkgerrors.IsOutdatedResourceError(outdatedError))
	assert.False(t, pkgerrors.IsBadRequestError(outdatedError))
	assert.Equal(t, "outdated", outdatedError.Error())
}
