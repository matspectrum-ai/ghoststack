.PHONY: build test lint fmt vet clean ci help

BINARY_NAME=ghost
VERSION?=$(shell cat VERSION 2>/dev/null || echo 0.1.0)
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

help:
	@echo "GhostStack Makefile"
	@echo ""
	@echo "Uso:"
	@echo "  make build        Compila o binário localmente"
	@echo "  make test         Executa testes unitários e de integração"
	@echo "  make lint        Executa linter e formatação"
	@echo "  make fmt          Aplica formatação automática"
	@echo "  make vet          Executa go vet"
	@echo "  make clean        Remove artefatos de build"
	@echo "  make ci           Pipeline local equivalente ao CI"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/ghost

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/ dist/ tmp/

ci: fmt vet lint test build
	@echo "CI local concluído com sucesso"
