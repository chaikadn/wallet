package storage

import "errors"

type Storage interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}

var (
	ErrWalletNotFound = errors.New("wallet not found")
)
