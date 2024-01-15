package models

import "errors"

var (
	ErrInsufficientFunds           = errors.New("insufficient funds")
	ErrTransactionAlreadyProcessed = errors.New("transaction already processed")
	ErrTransactionNotFound         = errors.New("transaction not found")
	ErrUserWalletNotFound          = errors.New("user wallet not found")

	ErrAmountNegative     = errors.New("amount should not be negative")
	ErrEmptyUserID        = errors.New("user id must not be empty")
	ErrEmptyTransactionID = errors.New("transaction id must not be empty")
)
