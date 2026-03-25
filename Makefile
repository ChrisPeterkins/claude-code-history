.PHONY: build install

build:
	go build -o claude-history .

install: build
	sudo cp claude-history /usr/local/bin/claude-history
