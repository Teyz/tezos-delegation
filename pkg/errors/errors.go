package errors

import (
	"github.com/pkg/errors"
)

type ErrorWithKey struct {
	Key string
}

func (e *ErrorWithKey) Error() string {
	return e.Key
}

// NewNotFoundError return a new NotFoundError
func NewNotFoundError(key string) error {
	return errors.WithStack(&NotFoundError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// NotFoundError is used when we cannot find a specified resource
type NotFoundError struct {
	ErrorWithKey
}

// IsNotFoundError verify if an error is a NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := errors.Cause(err).(*NotFoundError)

	return ok
}

// NewBadRequestError return a new BadRequestError
func NewBadRequestError(key string) error {
	return errors.WithStack(&BadRequestError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// BadRequestError is used when the given parameters does not match requirements
type BadRequestError struct {
	ErrorWithKey
}

// IsBadRequestError verify if an error is a BadRequestError
func IsBadRequestError(err error) bool {
	_, ok := errors.Cause(err).(*BadRequestError)

	return ok
}

// NewExpiredResourceError return a new ExpiredResourceError
func NewExpiredResourceError(key string) error {
	return errors.WithStack(&ExpiredResourceError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// ExpiresResourceError is used when the given resource has expired
type ExpiredResourceError struct {
	ErrorWithKey
}

// IsExpiresResourceError verify if an error is a ExpiredResourceError
func IsExpiredResourceError(err error) bool {
	_, ok := errors.Cause(err).(*ExpiredResourceError)

	return ok
}

// NewInternalServerError return a new InternalServerError
func NewInternalServerError(key string) error {
	return errors.WithStack(&InternalServerError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// InternalServerError is used when an error unexpected appears
type InternalServerError struct {
	ErrorWithKey
}

// IsInternalServerError verify if an error is a InternalServerError
func IsInternalServerError(err error) bool {
	_, ok := errors.Cause(err).(*InternalServerError)

	return ok
}

// NewUnauthorizedError return a new UnauthorizedError
func NewUnauthorizedError(key string, subjectAndMessage ...string) error {
	return errors.WithStack(&UnauthorizedError{
		ErrorWithKey: ErrorWithKey{key},
	})
}

// UnauthorizedError is used when action is not authorized
type UnauthorizedError struct {
	ErrorWithKey
}

// IsUnauthorizedError verify if an error is a UnauthorizedError
func IsUnauthorizedError(err error) bool {
	_, ok := errors.Cause(err).(*UnauthorizedError)

	return ok
}

// NewResourceAlreadyExist return a new ResourceAlreadyExist
func NewResourceAlreadyCreatedError(key string) error {
	return errors.WithStack(&ResourceAlreadyCreatedError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// ResourceAlreadyCreatedError is used when a resource already exist and could not be created another time
type ResourceAlreadyCreatedError struct {
	ErrorWithKey
}

// IsResourceAlreadyCreatedError verify if an error is a ResourceAlreadyCreatedError
func IsResourceAlreadyCreatedError(err error) bool {
	_, ok := errors.Cause(err).(*ResourceAlreadyCreatedError)

	return ok
}

// NewOutdatedResourceError return a new OutdatedResourceError
func NewOutdatedResourceError(key string) error {
	return errors.WithStack(&OutdatedResourceError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// ResourceAlreadyCreatedError is used when a resource already exist and could not be created another time
type OutdatedResourceError struct {
	ErrorWithKey
}

// IsResourceAlreadyCreatedError verify if an error is a ResourceAlreadyCreatedError
func IsOutdatedResourceError(err error) bool {
	_, ok := errors.Cause(err).(*OutdatedResourceError)

	return ok
}
