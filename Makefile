# Build automation and commands

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

# Binary names
BINARY_NAME=algo-scales
SERVER_BINARY_NAME=algo-scales-server

# Main package paths
MAIN_PATH=./
SERVER_PATH=./server

# Targets
.PHONY: all build clean test test-coverage fmt lint run server install vet

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	$(GOBUILD) -o $(SERVER_BINARY_NAME) -v $(SERVER_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(SERVER_BINARY_NAME)
	rm -f $(BINARY_NAME).test
	rm -f $(SERVER_BINARY_NAME).test
	rm -f coverage.out

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

fmt:
	$(GOFMT) ./...

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

server:
	$(GOBUILD) -o $(SERVER_BINARY_NAME) -v $(SERVER_PATH)
	./$(SERVER_BINARY_NAME)

install:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	sudo mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Cross-compile targets
.PHONY: build-linux build-windows build-darwin

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 -v $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(SERVER_BINARY_NAME)-linux-amd64 -v $(SERVER_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe -v $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(SERVER_BINARY_NAME)-windows-amd64.exe -v $(SERVER_PATH)

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 -v $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(SERVER_BINARY_NAME)-darwin-amd64 -v $(SERVER_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64 -v $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(SERVER_BINARY_NAME)-darwin-arm64 -v $(SERVER_PATH)

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin
