.DEFAULT_GOAL := dev
.PHONY: all dev server web build

BINARY=frontend
CMD=./server/cmd/main.go

dev: all
all:
	@$(MAKE) -j2 server web
server:
	find . -type f \( -name '*.go' -o -name '*.html' \) | entr -r go run $(CMD)
web:
	pnpm run dev:css
build:
	go build -o bin/$(BINARY) $(CMD)
