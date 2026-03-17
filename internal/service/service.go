package service

import "errors"

var (
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrServiceUnavailable = errors.New("service unavailable")
)
