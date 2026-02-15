package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server
	GraphQL
}

type Server struct {
	Port int `env:"SERVER_PORT" envDefault:"3000"`
}

type GraphQL struct {
	Url string `env:"GRAPHQL_URL" envDefault:"http://localhost:8080/query"`
}

func NewConfig() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Err(err).Msg("Failed to decode env")
	}

	return &c
}
