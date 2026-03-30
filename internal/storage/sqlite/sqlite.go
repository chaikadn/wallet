package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"wallet/internal/storage"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	// sqlite не может нормально работать при конкурентном доступе поэтому пока ограничим пул соединений до одного
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wallet(
			id TEXT PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0
				CHECK (balance >= 0)
		)`)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetBalance(walletID string) (int, error) {
	var balance int
	err := s.db.QueryRow(`SELECT balance FROM wallet WHERE id = ?`, walletID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, storage.ErrWalletNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	return balance, nil
}

func (s *Storage) Deposit(walletID string, amount int) error {
	res, err := s.db.Exec(`UPDATE wallet SET balance = balance + ? WHERE id = ?`, amount, walletID)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	rowsAffected, resErr := res.RowsAffected()
	if resErr != nil {
		return fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, resErr)
	}
	if rowsAffected == 0 {
		return storage.ErrWalletNotFound
	}

	return nil
}

func (s *Storage) Withdraw(walletID string, amount int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}
	defer tx.Rollback()

	var balance int
	err = s.db.QueryRow(`SELECT balance FROM wallet WHERE id = ?`, walletID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrWalletNotFound
	}

	if balance < amount {
		return storage.ErrBalanceCondition
	}

	_, err = s.db.Exec(`UPDATE wallet SET balance = balance - ? WHERE id = ?`, amount, walletID)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
