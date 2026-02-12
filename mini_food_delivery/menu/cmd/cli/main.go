package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/menu/db"
	"github.com/defvova/go_projects/mini_food_delivery/menu/internal/config"
	"github.com/defvova/go_projects/mini_food_delivery/menu/internal/seed"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(consoleWriter).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	c := config.NewConfig()

	dbCreds := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbCreds)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("transaction error")
		os.Exit(1)
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)

	if err := seed.SeedAll(ctx, q); err != nil {
		log.Fatal().Err(err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("menu seeding completed")
}
