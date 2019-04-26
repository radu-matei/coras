BINDIR          := $(CURDIR)/bin

.PHONY: build
build:
	go build -o bin/coras ./cmd/...

.PHONY: bootstrap
bootstrap:
	dep ensure -v
