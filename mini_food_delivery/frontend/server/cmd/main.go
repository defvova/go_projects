package main

import (
	"net/http"
	"strconv"

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
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))

	r.Get("/", routes.HomeHandler)

	p := ":" + strconv.Itoa(conf.Server.Port)
	log.Info().Msgf("Server starting on http://localhost%v", p)
	err := http.ListenAndServe(p, r)
	if err != nil {
		log.Fatal().Err(err)
	}
}
