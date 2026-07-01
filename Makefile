SHELL=bash

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: audit
audit:
	dis-vulncheck

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	exit
