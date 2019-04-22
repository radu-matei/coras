BINDIR          := $(CURDIR)/bin

.PHONY: build
build:
	go build -o bin/coras ./...

.PHONY: bootstrap
bootstrap:
	dep ensure -v
