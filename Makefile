VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build install

build:
	go build -ldflags "-X main.version=$(VERSION)" -o claude-history .

install: build
	sudo cp claude-history /usr/local/bin/claude-history
