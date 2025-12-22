SHELL := /bin/sh

GO ?= go
CMD ?= ./cmd/mysql-client-gui
BIN_DIR ?= output
BINARY ?= mysql-client-gui
CGO_ENABLED ?= 0

.PHONY: all build run test clean

all: build

build:
	@mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -o "$(BIN_DIR)/$(BINARY)" $(CMD)

run:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) run $(CMD)

test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./...

clean:
	rm -rf "$(BIN_DIR)"
