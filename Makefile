#
# github.com/MichaelDarr/ahab
#

BIN := ahab
DESTDIR :=
GO ?= go
PREFIX := /usr/local
VERSION = $(shell cat VERSION)

EXTRA_GOFLAGS ?=
LDFLAGS := $(LDFLAGS) -X "github.com/MichaelDarr/ahab/internal.Version=$(VERSION)"

.PHONY: default
default: $(BIN)

.PHONY: self
self: ## use ahab to build itself
	$(BIN) cmd make

.PHONY: build
build: $(BIN)

.PHONY: $(BIN)
$(BIN): ## build
	$(GO) build $(GOFLAGS) -ldflags '-s -w $(LDFLAGS)' $(EXTRA_GOFLAGS) -o $@

.PHONY: test
test: ## use ahab to test itself
	cd test; \
	 $(BIN) cmd make containertest

.PHONY: containertest
containertest: ## must be run inside container set up for test suite
	$(GO) test -v -ldflags '-s -w $(LDFLAGS)' github.com/MichaelDarr/ahab/internal

.PHONY: install
install:
	install -Dm755 ${BIN} $(DESTDIR)$(PREFIX)/bin/${BIN}

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/${BIN}
