#
# github.com/MichaelDarr/docker-config
#

BIN := dcfg
GO ?= go
VERSION = $(shell cat VERSION)

GOFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: default
default: $(BIN)

.PHONY: build
build: $(BIN)

$(BIN): ## build docker-config as dcfg
	$(GO) build $(GOFLAGS) -o $(BIN) .

.PHONY: run
run: ## build and run docker-config
	$(GO) run $(GOFLAGS) .
