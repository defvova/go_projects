.DEFAULT_GOAL := dev
.PHONY: run dev

BINARY=webscraper
CMD=./cmd/webscraper/main.go

dev:
	find . -name '*.go' | entr -r go run $(CMD)
build:
	go build -o bin/$(BINARY) $(CMD)
