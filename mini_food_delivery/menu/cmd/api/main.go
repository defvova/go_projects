package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/menu/db"
	"github.com/defvova/go_projects/mini_food_delivery/menu/internal/config"
	"github.com/defvova/go_projects/mini_food_delivery/menu/internal/grpcserver"
	"github.com/defvova/go_projects/mini_food_delivery/menu/internal/observability"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pgxzerolog "github.com/jackc/pgx-zerolog"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := log.Output(consoleWriter).With().Timestamp().Logger()
	log.Logger = logger

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	c := config.NewConfig()
	ctx := context.Background()

	if c.Observability.OtelEnabled {
		shutdown, err := observability.InitOtel(
			ctx,
			c.Observability.ServiceName,
			c.Observability.TempoUrl,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("otel init failed")
		}
		defer shutdown(ctx)
	}

	dbCreds := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	cfg, err := pgxpool.ParseConfig(dbCreds)
	if err != nil {
		panic(err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzerolog.NewLogger(logger),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	queries := db.New(pool)

	if c.Observability.PrometheusEnabled {
		metricPort := ":" + strconv.Itoa(c.Observability.OtelPort)
		go func() {
			log.Info().Msgf("metrics listening on %v", metricPort)
			if err := http.ListenAndServe(metricPort, promhttp.Handler()); err != nil {
				log.Fatal().Err(err).Msg("metrics server failed")
			}
		}()
	}

	grpcServer := grpcserver.NewGRPCServer(
		queries,
		grpcserver.GRPCOptions{
			OtelEnabled:       c.Observability.OtelEnabled,
			PrometheusEnabled: c.Observability.PrometheusEnabled,
		},
	)
	lis, err := grpcServer.Listen(":" + strconv.Itoa(c.GRPCServer.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	go func() {
		log.Info().Msgf("gRPC listening on %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("grpc serve failed")
		}
	}()

	waitForShutdown(grpcServer.Server)
}

func waitForShutdown(server *grpc.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("shutting down gRPC server...")
	server.GracefulStop()
}
