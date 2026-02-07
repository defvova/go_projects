package main

import (
	"context"
	"fmt"
	"mini_food_delivery/menu/db"
	"mini_food_delivery/menu/internal/config"
	"mini_food_delivery/menu/internal/grpcserver"
	"mini_food_delivery/menu/internal/otel"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	c := config.NewConfig()
	appCtx := context.Background()

	shutdown, err := otel.Init(
		appCtx,
		"menu-service",
		"localhost:4317",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("otel init failed")
	}
	defer shutdown(appCtx)

	s := grpcserver.NewGRPCServer(":" + strconv.Itoa(c.GRPCServer.Port))
	dbCreds := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbCreds)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	server, lis, err := s.ServeListener(db.New(conn))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run server")
	}

	go func() {
		log.Info().Msgf("gRPC listening on %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatal().Err(err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *grpc.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("shutting down gRPC server...")
	server.GracefulStop()
}
