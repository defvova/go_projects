package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/resolver"
	"github.com/defvova/go_projects/mini_food_delivery/gateway/internal/config"
	"github.com/defvova/go_projects/mini_food_delivery/gateway/internal/menugrpc"
	"github.com/defvova/go_projects/mini_food_delivery/gateway/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	c := config.NewConfig()
	r := chi.NewRouter()
	port := ":" + strconv.Itoa(c.Server.Port)
	s := &http.Server{
		Addr:         port,
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

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{
		MenuService: services.MenuService{
			MenuClient: menuGrpc.Menu,
			// MenuItemClient: menuGrpc.MenuItem,
		},
		CategoryService: services.CategoryService{
			CategoryClient: menuGrpc.Category,
		},
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

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

	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure!")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
