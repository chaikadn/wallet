package wallet

import (
	"fmt"
	"wallet/internal/storage"
)

type WalletService struct {
	st storage.Storage
}

func (ws *WalletService) GetBalance(walletID string) (int, error) {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	balance, err := ws.st.GetBalance(walletID)
	if err != nil {
		return 0, fmt.Errorf("service error: get balance: %v", err)
	}
	return balance, nil
}

func (ws *WalletService) Deposit(walletID string, amount int) error {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	err := ws.st.Deposit(walletID, amount)
	if err != nil {
		return fmt.Errorf("service error: deposit: %v", err)
	}
	return nil
}

func (ws *WalletService) Withdraw(walletID string, amount int) error {
	// валидировать uuid на бизнес-правила (есть ли такой кошелек)

	err := ws.st.Withdraw(walletID, amount)
	if err != nil {
		return fmt.Errorf("service error: withdraw: %v", err)
	}
	return nil
}
