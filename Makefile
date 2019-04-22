BINDIR          := $(CURDIR)/bin

.PHONY: build
build:
	go build -o bin/coras ./...
