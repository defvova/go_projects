package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"shortener/config"
	"shortener/internal/router"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	a := config.NewApp()
	dbCreds := fmt.Sprintf(fmtDBString, a.ConfDb.Host, a.ConfDb.Username, a.ConfDb.Password, a.ConfDb.DBName, a.ConfDb.Port)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbCreds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	handler := router.InitHandlers(mux, conn)

	log.Println("Server starting on" + a.Addr)
	log.Fatal(http.ListenAndServe(a.Addr, handler))
}
