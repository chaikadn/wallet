package pg

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wallet/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(uri string) (*Storage, error) {
	// pool, err := pgxpool.New(context.TODO(), uri)

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	pool, err := pgxpool.NewWithConfig(context.TODO(), config)

	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(context.TODO(),
		`CREATE TABLE IF NOT EXISTS wallet(
			uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			balance BIGINT NOT NULL DEFAULT 0,
			
			CONSTRAINT balance_not_negative CHECK (balance >= 0)
		)`)
	if err != nil {
		return nil, err
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) GetBalance(ctx context.Context, walletID string) (int, error) {
	var balance int
	err := s.pool.QueryRow(ctx, `SELECT balance FROM wallet WHERE uuid = $1`, walletID).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, storage.ErrWalletNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("%w: %w", storage.ErrDatabaseUnavailable, err)
	}

	return balance, nil
}

func (s *Storage) Deposit(ctx context.Context, walletID string, amount int) error {
	res, err := s.pool.Exec(ctx, `UPDATE wallet SET balance = balance + $1 WHERE uuid = $2`, amount, walletID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return storage.ErrWalletNotFound
	}

	return nil
}

// TODO: лучше сделать атомарный UPDATE и ловить rows affected
func (s *Storage) Withdraw(ctx context.Context, walletID string, amount int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var balance int
	err = tx.QueryRow(ctx, `SELECT balance FROM wallet WHERE uuid = $1`, walletID).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrWalletNotFound
	}
	if err != nil {
		return err
	}

	if balance < amount {
		return storage.ErrBalanceCondition
	}

	_, err = tx.Exec(ctx, `UPDATE wallet SET balance = balance - $1 WHERE uuid = $2`, amount, walletID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return storage.ErrBalanceCondition
		}
		return err
	}

	return tx.Commit(ctx)
}

// TODO
func (s *Storage) Close() error {
	s.pool.Close()
	return nil
}
