package errs

import "errors"

var (
	ErrRecordNotFound = errors.New("This record isn't found")
	ErrInternal       = errors.New("internal system error")

	ErrUnauthorized = errors.New("Unauthorized")
)
