package wallet

import (
	"errors"
	"wallet/internal/service"
	"wallet/internal/storage"

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
	balance, err := w.walletStorage.GetBalance(walletID)
	if errors.Is(err, storage.ErrWalletNotFound) {
		return 0, service.ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("get balance error", zap.Error(err))
		return 0, service.ErrServiceUnavailable
	}
	return balance, nil
}

func (w *Wallet) Deposit(walletID string, amount int) error {
	if amount <= 0 {
		return service.ErrInvalidAmount
	}

	err := w.walletStorage.Deposit(walletID, amount)
	if errors.Is(err, storage.ErrWalletNotFound) {
		return service.ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("deposit error", zap.Error(err))
		return service.ErrServiceUnavailable
	}
	return nil
}

func (w *Wallet) Withdraw(walletID string, amount int) error {
	// гонка данных между GetBalance и Withdraw (между этими проверками баланс может измениться)
	// решается либо атомарной операцией в самом репозитории (UPDATE ... WHERE balance >= amount), либо транзакцией с SELECT FOR UPDATE
	// лучше использовать транзакцию

	if amount <= 0 {
		return service.ErrInvalidAmount
	}

	balance, err := w.walletStorage.GetBalance(walletID)

	if errors.Is(err, storage.ErrWalletNotFound) {
		return service.ErrWalletNotFound
	}
	if err != nil {
		w.log.Error("withdraw error", zap.Error(err))
		return service.ErrServiceUnavailable
	}
	if amount > balance {
		return service.ErrInsufficientFunds
	}

	err = w.walletStorage.Withdraw(walletID, amount)
	if err != nil {
		w.log.Error("withdraw error", zap.Error(err))
		return service.ErrServiceUnavailable
	}
	return nil
}
