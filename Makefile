#
# github.com/MichaelDarr/ahab
#

BIN := ahab
DESTDIR :=
GO ?= go
PREFIX := /usr/local
VERSION = $(shell cat VERSION)

GOFLAGS := -mod=vendor
EXTRA_GOFLAGS ?=
LDFLAGS := $(LDFLAGS) -X "github.com/MichaelDarr/ahab/internal.Version=$(VERSION)"

.PHONY: default
default: $(BIN)

.PHONY: self
self: ## use ahab to build itself
	$(BIN) exec make

.PHONY: build
build: $(BIN)

.PHONY: $(BIN)
$(BIN): ## build
	$(GO) build $(GOFLAGS) -ldflags '-s -w $(LDFLAGS)' $(EXTRA_GOFLAGS) -o $@

.PHONY: test
test: ## use ahab to test itself
	cd test; \
	 $(BIN) exec make containertest

.PHONY: containertest
containertest: ## must be run inside container set up for test suite
	$(GO) test $(GOFLAGS) -ldflags '-s -w $(LDFLAGS)' github.com/MichaelDarr/ahab/internal

.PHONY: coverage
coverage: ## use ahab to test itself and generate a coverage report
	cd test; \
	 $(BIN) exec make containercoverage

.PHONY: containercoverage
containercoverage: ## also run inside container, with verbose output and a coverage report
	$(GO) test $(GOFLAGS) -v -coverprofile cp.out -ldflags '-s -w $(LDFLAGS)' github.com/MichaelDarr/ahab/internal

.PHONY: install
install:
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 0755 $(BIN) $(DESTDIR)$(PREFIX)/bin/$(BIN)

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BIN)
