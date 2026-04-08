package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	DatabaseDSN string `env:"DATABASE_DSN,required"`
	HTTPServer  `envPrefix:"SERVER_"`
}

type HTTPServer struct {
	Address     string        `env:"ADDRESS" envDefault:"0.0.0.0:8080"`
	Timeout     time.Duration `env:"TIMEOUT" envDefault:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	// user password для админки
}

func MustLoad() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	return cfg
}
