package wallet

import (
	"fmt"

	"go.uber.org/zap"
)

type walletStorage interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}

type Wallet struct {
	walletStorage walletStorage
	log           *zap.Logger
}

func NewWallet(walletStorage walletStorage, logger *zap.Logger) *Wallet {
	return &Wallet{walletStorage: walletStorage, log: logger}
}

func (w *Wallet) GetBalance(walletID string) (int, error) {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	balance, err := w.walletStorage.GetBalance(walletID)
	if err != nil {
		return 0, fmt.Errorf("service error: get balance: %v", err)
	}
	return balance, nil
}

func (w *Wallet) Deposit(walletID string, amount int) error {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	err := w.walletStorage.Deposit(walletID, amount)
	if err != nil {
		return fmt.Errorf("service error: deposit: %v", err)
	}
	return nil
}

func (w *Wallet) Withdraw(walletID string, amount int) error {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	err := w.walletStorage.Withdraw(walletID, amount)
	if err != nil {
		return fmt.Errorf("service error: withdraw: %v", err)
	}
	return nil
}
