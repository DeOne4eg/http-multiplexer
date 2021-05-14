.SILENT:
.PHONY: build
build:
	go build -o .build/ ./cmd/multiplexer

.PHONY: run
run:
	go run ./cmd/multiplexer.go

.PHONY: test
test:
	go test -coverprofile=cover.out -v ./...

.PHONY: test.coverage
test.coverage:
	go tool cover -func=cover.out

.PHONY: lint
lint:
	golangci-lint run ./...

.DEFAULT_GOAL := build