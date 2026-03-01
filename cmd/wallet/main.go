package main

import (
	"log"
	"net/http"
	"wallet/internal/handler"
	"wallet/internal/service/wallet"
	"wallet/internal/storage/pg"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	logger, err := newLogger("debug")
	if err != nil {
		log.Fatal("failed to initialize logger", err)
	}
	defer logger.Sync()

	storage := pg.NewPGStorage(logger)
	service := wallet.NewWalletService(storage, logger)
	handler := handler.NewHandler(service, logger)

	router := chi.NewMux()
	handler.RegisterRoutes(router)

	logger.Info("starting server", zap.String("address", ":8080"))
	http.ListenAndServe(":8080", router)
}

func newLogger(logLevel string) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = level
	cfg.DisableStacktrace = true

	return cfg.Build()
}
