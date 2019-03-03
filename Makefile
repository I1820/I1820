VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

## install: Install missing dependencies.
install: go-get

## build: Build the binary.
build: go-build

## Lint: Lint all source codes
lint: go-lint

go-build:
	@echo "  >  Building binary..."
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	go get

go-lint:
	@echo "  >  Linting with 'golangci-lint'..."
	golangci-lint run --enable-all ./...
