package sqlite

import (
	"database/sql"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string, logger *zap.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetBalance(walletID string) (int, error) {
	var balance int
	err := s.db.QueryRow("SELECT balance FROM wallet WHERE id = ?", walletID).Scan(&balance)
	if err != nil {
		// handle sql errors
		return 0, err
	}
	return balance, nil
}

func (s *Storage) Deposit(walletID string, amount int) error {
	_, err := s.db.Exec("UPDATE wallet SET balance = balance + ? WHERE id = ?", amount, walletID)
	if err != nil {
		// handle sql errors
		return err
	}
	return nil
}

func (s *Storage) Withdraw(walletID string, amount int) error {
	_, err := s.db.Exec("UPDATE wallet SET balance = balance - ? WHERE id = ?", amount, walletID)
	if err != nil {
		// handle sql errors
		return err
	}
	return nil
}

// _ "github.com/mattn/go-sqlite3" использует CGO
