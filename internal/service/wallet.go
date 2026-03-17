package service

import (
	"errors"

	"wallet/internal/storage"

	"go.uber.org/zap"
)

type WalletStorage interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}

type Wallet struct {
	walletStorage WalletStorage
	log           *zap.Logger
}

func NewWallet(walletStorage WalletStorage, logger *zap.Logger) *Wallet {
	return &Wallet{walletStorage: walletStorage, log: logger}
}

func (w *Wallet) GetBalance(walletID string) (int, error) {
	balance, err := w.walletStorage.GetBalance(walletID)
	if errors.Is(err, storage.ErrWalletNotFound) {
		return 0, ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("get balance error", zap.Error(err))
		return 0, ErrServiceUnavailable
	}
	return balance, nil
}

func (w *Wallet) Deposit(walletID string, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := w.walletStorage.Deposit(walletID, amount)
	if errors.Is(err, storage.ErrWalletNotFound) {
		return ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("deposit error", zap.Error(err))
		return ErrServiceUnavailable
	}
	return nil
}

func (w *Wallet) Withdraw(walletID string, amount int) error {
	// гонка данных между GetBalance и Withdraw (между этими проверками баланс может измениться)
	// решается либо атомарной операцией в самом репозитории (UPDATE ... WHERE balance >= amount), либо транзакцией с SELECT FOR UPDATE
	// лучше использовать транзакцию

	if amount <= 0 {
		return ErrInvalidAmount
	}

	balance, err := w.walletStorage.GetBalance(walletID)

	if errors.Is(err, storage.ErrWalletNotFound) {
		return ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("withdraw error", zap.Error(err))
		return ErrServiceUnavailable
	}
	if amount > balance {
		return ErrInsufficientFunds
	}

	err = w.walletStorage.Withdraw(walletID, amount)
	if err != nil {
		w.log.Error("withdraw error", zap.Error(err))
		return ErrServiceUnavailable
	}
	return nil
}
