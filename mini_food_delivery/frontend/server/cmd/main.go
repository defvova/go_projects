package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/defvova/go_projects/mini_food_delivery/frontend/internal/config"
	"github.com/defvova/go_projects/mini_food_delivery/frontend/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	conf := config.NewConfig()
	port := ":" + strconv.Itoa(conf.Server.Port)
	r := chi.NewRouter()
	s := &http.Server{
		Addr:    port,
		Handler: r,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))

	r.Get("/", routes.HomeHandler)

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msgf("Shutting down server %v", port)

		ctx, cancel := context.WithTimeout(context.Background(), 60)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failure")
		}

		close(closed)
	}()

	log.Info().Msgf("Server listening on http://localhost%v", port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure!")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
