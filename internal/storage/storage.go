package storage

import "errors"

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrBalanceCondition    = errors.New("balance condition failed")
	ErrDatabaseUnavailable = errors.New("database unavailable")
)
