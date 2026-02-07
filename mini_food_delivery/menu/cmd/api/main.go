package main

import (
	"context"
	"fmt"
	"mini_food_delivery/menu/db"
	"mini_food_delivery/menu/internal/config"
	"mini_food_delivery/menu/internal/grpcserver"
	"mini_food_delivery/menu/internal/observability"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(consoleWriter).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	c := config.NewConfig()
	appCtx := context.Background()

	if c.Observability.OtelEnabled {
		shutdown, err := observability.InitOtel(
			appCtx,
			c.Observability.ServiceName,
			c.Observability.TempoUrl,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("otel init failed")
		}
		defer shutdown(appCtx)
	}

	dbCreds := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbCreds)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

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
		db.New(conn),
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
