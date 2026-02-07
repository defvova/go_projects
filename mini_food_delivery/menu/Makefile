.DEFAULT_GOAL := server
.PHONY: server

BINARY=menu
API_CMD=./cmd/api/main.go
MIGRATION_CMD=./cmd/migrate/main.go

PROTO_DIR=internal/protobuf
PROTO_SRC=$(wildcard $(PROTO_DIR)/**/**/*.proto)
GO_OUT=./pkg

server:
	find . -type f \( -name '*.go' -o -name '*.html' \) | entr -r go run $(API_CMD)
db-seed:
	go run ./cmd/cli/main.go
db-g-%:
	goose create $* sql
db-%:
	go run $(MIGRATION_CMD) $*
proto-gen:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)
