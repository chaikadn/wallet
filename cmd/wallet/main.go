package main

import (
	"log"
	"net/http"
	"wallet/internal/handler"
	"wallet/internal/service/wallet"
	"wallet/internal/storage/pg"

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

	logger.Info("starting server", zap.String("address", ":8080"))
	http.ListenAndServe(":8080", handler.RegisterRoutes())

	// TODO: grasefull shutdown
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

// router := chi.NewMux()
// handler.RegisterRoutes(router)

// Нет таймаутов сервера. Каждое keep-alive соединение будет работать в отдельной горутине, а горутина никогда не завершится.
// надо делать свой Server и указывать таймауты. для данного проекта пока можно без таймаутов,
// т.к. все клиенты (например curl) "одноразовые" и соединения закрываются после завершения их работы.
