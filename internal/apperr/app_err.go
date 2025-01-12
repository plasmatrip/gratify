package apperr

import "errors"

var (
	ErrBadLogin          = errors.New("bad login or password")
	ErrLoginAlreadyTaken = errors.New("login already taken")

	ErrZeroRowInsert = errors.New("zero rows inserted")
)
