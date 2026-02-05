package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server APIServer
	MenuGRPCServer
}

type MenuGRPCServer struct {
	Url string `env:"MENU_GRPC_URL" envDefault:"localhost:5001"`
}

type APIServer struct {
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ" envDefault:"30s"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE" envDefault:"30s"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE" envDefault:"10s"`
}

func NewConfig() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode env")
	}

	return &c
}
