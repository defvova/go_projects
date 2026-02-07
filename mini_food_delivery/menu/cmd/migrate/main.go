package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"mini_food_delivery/menu/internal/config"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	dialect     = "pgx"
	fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "path to migrations directory")

	//go:embed migrations/*.sql
	embedMigrations embed.FS
)

func main() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(consoleWriter).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	command := args[0]

	c := config.NewDb()
	dbString := fmt.Sprintf(fmtDBString, c.Host, c.Username, c.Password, c.DBName, c.Port)

	db, err := goose.OpenDBWithDriver(dialect, dbString)
	if err != nil {
		log.Fatal().Err(err).Msg("open db")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("close db")
		}
	}()

	goose.SetBaseFS(embedMigrations)

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir, args[1:]...); err != nil {
		log.Fatal().Msgf("migrate %v: %v", command, err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
