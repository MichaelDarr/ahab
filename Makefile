#
# github.com/MichaelDarr/ahab
#

BIN := ahab
ENTRYPOINT := ahab.go
GO ?= go
VERSION = $(shell cat VERSION)

GOFLAGS := -ldflags "-X github.com/MichaelDarr/ahab/internal.Version=$(VERSION)"

.PHONY: default
default: $(BIN)

.PHONY: build
build: $(BIN)

.PHONY: $(BIN)
$(BIN): ## build
	$(GO) build $(GOFLAGS) $(ENTRYPOINT)

.PHONY: run
run: ## build and run
	$(GO) run $(GOFLAGS) $(ENTRYPOINT)
