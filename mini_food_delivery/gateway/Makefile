.DEFAULT_GOAL := dev
.PHONY: run dev

BINARY=mini_food_delivery
CMD=./cmd/api/main.go

dev:
	find . -name '*.go' | entr -r go run $(CMD)
build:
	go build -o bin/$(BINARY) $(CMD)
gqlgen:
	go run github.com/99designs/gqlgen generate
