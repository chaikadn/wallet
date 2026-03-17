package main

import (
	defaultLog "log"
	"net/http"
	"wallet/internal/config"
	"wallet/internal/handler"
	"wallet/internal/service"
	"wallet/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	log, err := newLogger(cfg.LogLevel)
	if err != nil {
		defaultLog.Fatal("failed to initialize logger", err)
	}
	defer log.Sync()

	log.Info("starting wallet", zap.String("log-level", cfg.LogLevel))

	storage, err := sqlite.New(cfg.StoragePath, log)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer storage.Close()

	log.Debug("successfully connected to db", zap.String("path", cfg.StoragePath))

	wallet := service.NewWallet(storage, log)

	router := chi.NewRouter()

	handler := handler.NewHandler(wallet, log)
	handler.RegisterRoutes(router)

	srv := newServer(cfg, router)

	log.Info("starting server", zap.String("address", cfg.HTTPServer.Address))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	// TODO: grasefull shutdown
}

// можно заменить на slog, он в стандартной библиотеке с go 1.21
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

func newServer(cfg *config.Config, handler http.Handler) http.Server {
	return http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}
