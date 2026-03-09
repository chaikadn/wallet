package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	StoragePath string `env:"STORAGE_PATH,required"`
	HTTPServer  `envPrefix:"SERVER_"`
}

type HTTPServer struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	Timeout     string `env:"TIMEOUT" envDefault:"4s"`
	IdleTimeout string `env:"IDLE_TIMEOUT" envDefault:"60s"`
	// user password для админки
}

func MustLoad() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	return cfg
}
