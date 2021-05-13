.PHONY: build
build:
	go build -o .build/ ./cmd/multiplexer

.PHONE: run
run:
	go run ./cmd/multiplexer.go

.PHONY: lint
lint:
	golangci-lint run ./...

.DEFAULT_GOAL := build