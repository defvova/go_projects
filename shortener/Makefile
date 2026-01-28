.DEFAULT_GOAL := dev
.PHONY: all server dev web build

BINARY=shortener
CMD=./cmd/shortener/main.go

all:
	@$(MAKE) -j2 server web
dev: all
server:
	find . -type f \( -name '*.go' -o -name '*.html' \) | entr -r go run $(CMD)
web:
	pnpm run dev:css
db-migrate:
	@set -a; . .env; set +a; \
	psql -U "$$DB_USER" -d "$$DB_NAME" -f schema.sql
build:
	go build -o bin/$(BINARY) $(CMD)
