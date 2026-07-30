.PHONY: build build-all test test-race test-coverage lint fmt vet clean ci fuzz benchmark help cover-html

BINARY_NAME=ghost
VERSION?=$(shell cat VERSION 2>/dev/null || echo 0.1.0)
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X github.com/ghoststack/ghoststack/internal/cli.Version=$(VERSION) -X github.com/ghoststack/ghoststack/internal/cli.Commit=$(COMMIT) -X github.com/ghoststack/ghoststack/internal/cli.BuildTime=$(BUILD_TIME)"

help:
	@echo "GhostStack Makefile"
	@echo ""
	@echo "Uso:"
	@echo "  make build         Compila o binário localmente (linux/amd64)"
	@echo "  make build-all     Compila para linux/amd64, linux/arm64, darwin/amd64"
	@echo "  make test          Executa testes unitários"
	@echo "  make test-race     Executa testes com race detector"
	@echo "  make test-coverage Executa testes com cobertura (coverage.out)"
	@echo "  make cover-html    Abre relatório de cobertura no navegador"
	@echo "  make lint          Executa linter e formatação"
	@echo "  make fmt           Aplica go fmt"
	@echo "  make vet           Executa go vet"
	@echo "  make fuzz          Executa fuzz tests (10s)"
	@echo "  make benchmark     Executa benchmarks"
	@echo "  make clean         Remove artefatos de build"
	@echo "  make ci            Pipeline local completo (fmt -> vet -> lint -> test -race -> build)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/ghost

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/ghost
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/ghost
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/ghost

test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic -count=1 ./...
	@go tool cover -func=coverage.out | tail -1

cover-html: test-coverage
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

fuzz:
	go test -fuzz=Fuzz -fuzztime=10s ./...

benchmark:
	go test -bench=. -benchmem ./...

docker:
	docker build -t ghoststack:latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		.

docker-run: docker
	docker compose up -d

docker-stop:
	docker compose down

clean:
	rm -rf bin/ dist/ coverage.out tmp/

ci: fmt vet test-race build
	@echo "CI local concluído com sucesso"
