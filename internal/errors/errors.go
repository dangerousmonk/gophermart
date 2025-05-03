package errors

import "errors"

var (
	//user
	ErrNoUserIDFound = errors.New("user: failed to resolve userID")
	ErrLoginExists   = errors.New("user: user with this login already exists")
	ErrWrongPassword = errors.New("user: password is wrong")
	ErrNoUserFound   = errors.New("user: user not found by login")

	//orders
	ErrWrongOrderNum          = errors.New("order: order number is invalid")
	ErrOrderExists            = errors.New("order: order already uploaded by user")
	ErrOrderExistsAnotherUser = errors.New("order: order already uploaded by different user")
	ErrNoOrders               = errors.New("order: no orders found for user")
	ErrUnknown                = errors.New("order: failed to upload order")

	//withdrawals
	ErrNoWithdrawals            = errors.New("withdrawals: no withdrawals found for user")
	ErrWithdrawalForOrderExists = errors.New("withdrawal: for this order already registered")
	ErrInsufficientBalance      = errors.New("withdrawal: insufficient funds on balance")

	ErrValidationFailed = errors.New("validation: failed to validate request body")
)
