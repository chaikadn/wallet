package pg

import "go.uber.org/zap"

type PGStorage struct {
	log *zap.Logger
}

func NewPGStorage(logger *zap.Logger) *PGStorage {
	return &PGStorage{log: logger}
}

func (ps *PGStorage) GetBalance(walletID string) (int, error) {

	ps.log.Info("mock get balance", zap.String("wallet_id", walletID))

	return 0, nil
}

func (ps *PGStorage) Deposit(walletID string, amount int) error {

	ps.log.Info("mock deposit", zap.String("wallet_id", walletID), zap.Int("amount", amount))

	return nil
}

func (ps *PGStorage) Withdraw(walletID string, amount int) error {

	ps.log.Info("mock withdraw", zap.String("wallet_id", walletID), zap.Int("amount", amount))

	return nil
}
