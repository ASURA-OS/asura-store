.PHONY: all build test clean run-ui

PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
DATADIR ?= $(PREFIX)/share/asura/store

all: build

build:
	cd cmd/asura-store && go build -o ../../bin/asura-store .

test:
	cd cmd/asura-store && go test -v ./...
	cd ui/lib && node --test store.test.mjs

run-ui:
	quickshell -p ui/store.qml

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 bin/asura-store $(DESTDIR)$(BINDIR)/asura-store
	install -d $(DESTDIR)$(DATADIR)
	cp -r ui $(DESTDIR)$(DATADIR)/
	cp -r catalog $(DESTDIR)$(DATADIR)/

clean:
	rm -rf bin/
