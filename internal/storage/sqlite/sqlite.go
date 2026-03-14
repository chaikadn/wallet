package sqlite

import (
	"database/sql"
	"errors"
	"sync"
	"wallet/internal/storage"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
	mx *sync.RWMutex
}

func New(dbPath string, logger *zap.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wallet(
			id TEXT PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0
				CHECK (balance >= 0)
		)`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db, mx: &sync.RWMutex{}}, nil
}

func (s *Storage) GetBalance(walletID string) (int, error) {
	// TODO: лучше транзакцию
	s.mx.RLock()
	defer s.mx.RUnlock()

	var balance int
	err := s.db.QueryRow(`SELECT balance FROM wallet WHERE id = ?`, walletID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, storage.ErrWalletNotFound
	}
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (s *Storage) Deposit(walletID string, amount int) error {
	// TODO: лучше транзакцию
	s.mx.Lock()
	defer s.mx.Unlock()

	res, err := s.db.Exec(`UPDATE wallet SET balance = balance + ? WHERE id = ?`, amount, walletID)
	if err != nil {
		return err
	}

	rowsAffected, resErr := res.RowsAffected()
	if resErr != nil {
		return resErr
	}
	if rowsAffected == 0 {
		return storage.ErrWalletNotFound
	}

	return nil
}

func (s *Storage) Withdraw(walletID string, amount int) error {
	// TODO: лучше транзакцию
	s.mx.Lock()
	defer s.mx.Unlock()

	// TODO: читать баланс до операции, чтобы проверять не ушел ли в минус баланс

	res, err := s.db.Exec(`UPDATE wallet SET balance = balance - ? WHERE id = ?`, amount, walletID)
	if err != nil {
		return err
	}

	rowsAffected, resErr := res.RowsAffected()
	if resErr != nil {
		return resErr
	}
	if rowsAffected == 0 {
		return storage.ErrWalletNotFound
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// _ "github.com/mattn/go-sqlite3" использует CGO
