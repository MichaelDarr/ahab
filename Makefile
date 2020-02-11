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

.PHONY: build
build: $(BIN)

.PHONY: $(BIN)
$(BIN): ## build
	$(GO) build $(GOFLAGS) -ldflags '-s -w $(LDFLAGS)' $(EXTRA_GOFLAGS) -o $@

.PHONY: install
install:
	install -Dm755 ${BIN} $(DESTDIR)$(PREFIX)/bin/${BIN}

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/${BIN}
