package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	GRPCServer GRPCServer
	DB         ConfigDB
	Observability
}

type ConfigDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	Debug    bool   `env:"DB_DEBUG,required"`
}

type GRPCServer struct {
	Port int `env:"GRPC_PORT" envDefault:"5001"`
}

type Observability struct {
	OtelEnabled       bool   `env:"OTEL_ENABLED" envDefault:"true"`
	PrometheusEnabled bool   `env:"pROMETHEUS_ENABLED" envDefault:"true"`
	OtelPort          int    `env:"OTEL_PORT" envDefault:"5090"`
	ServiceName       string `env:"SERVICE_NAME" envDefault:"menu-service"`
	TempoUrl          string `env:"TEMPO_URL" envDefault:"localhost:4317"`
}

func NewConfig() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode env")
	}

	return &c
}

func NewDb() *ConfigDB {
	var c ConfigDB
	if err := env.Parse(&c); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode env")
	}

	return &c
}
