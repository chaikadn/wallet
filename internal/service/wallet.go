package service

import (
	"context"
	"errors"

	"wallet/internal/storage"

	"go.uber.org/zap"
)

type WalletStorage interface {
	GetBalance(ctx context.Context, walletID string) (int, error)
	Deposit(ctx context.Context, walletID string, amount int) error
	Withdraw(ctx context.Context, walletID string, amount int) error
}

type Wallet struct {
	walletStorage WalletStorage
	log           *zap.Logger
}

func NewWallet(walletStorage WalletStorage, logger *zap.Logger) *Wallet {
	return &Wallet{walletStorage: walletStorage, log: logger}
}

func (w *Wallet) GetBalance(walletID string) (int, error) {
	balance, err := w.walletStorage.GetBalance(context.TODO(), walletID)
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

	err := w.walletStorage.Deposit(context.TODO(), walletID, amount)
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
	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := w.walletStorage.Withdraw(context.TODO(), walletID, amount)
	if errors.Is(err, storage.ErrBalanceCondition) {
		return ErrInsufficientFunds
	}
	if err != nil {
		w.log.Error("withdraw error", zap.Error(err))
		return ErrServiceUnavailable
	}
	return nil
}
