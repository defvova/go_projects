package main

import (
	"context"
	"fmt"
	services "mini_food_delivery/gateway/internal/api"
	"mini_food_delivery/gateway/internal/config"
	"mini_food_delivery/gateway/internal/menugrpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	c := config.NewConfig()
	r := chi.NewRouter()
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	menuGrpc, err := menugrpc.New(c.MenuGRPCServer.Url)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to menu grpc")
	}
	defer menuGrpc.Close()

	api := services.ApiService{
		MenuGrpc: menuGrpc,
		R:        r,
	}
	api.Init()

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msgf("Shutting down server %v", c.Server.Port)

		ctx, cancel := context.WithTimeout(context.Background(), 60)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failure")
		}

		close(closed)
	}()

	log.Info().Msgf("Server listening on port %v", c.Server.Port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure!")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
