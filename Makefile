GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=slog
BIN_DIR=/usr/local/bin

build:
	$(GOGET) -d ./...
	$(GOBUILD) -o $(BINARY_NAME) -v
	mv $(BINARY_NAME) $(BIN_DIR)

local:
	$(GOGET) -d ./...
	$(GOBUILD) -o $(BINARY_NAME) -v