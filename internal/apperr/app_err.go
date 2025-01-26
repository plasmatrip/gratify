package apperr

import "errors"

var (
	ErrBadLogin          = errors.New("bad login or password")
	ErrLoginAlreadyTaken = errors.New("login already taken")

	ErrZeroRowInsert                   = errors.New("zero rows inserted")
	ErrOrderAlreadyUploadedAnotherUser = errors.New("order number has already been uploaded by another user")
	ErrOrderIsNotRegisteredInAccrual   = errors.New("the order is not registered in the accrual system")

	ErrInternalServerAccrualError = errors.New("internal server accrual error")
)
