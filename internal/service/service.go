package service

type Service interface {
	GetBalance(walletID string) (int, error)
	Deposit(walletID string, amount int) error
	Withdraw(walletID string, amount int) error
}
