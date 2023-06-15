BINARY_NAME=dlc

all: help

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: mod
mod: ## Go mod things
	go mod tidy
	go mod vendor
	go mod download

build: clean ## Build app
	CGO=0 go build -o dist/${BINARY_NAME} cmd/main.go

clean: ## Clean stuff
	rm -rf dist/*

format: ## Format
	go fmt ./...

dev/install: build ## Install dlc locally
	rm -rf ~/.local/bin/dlc
	cp dist/dlc ~/.local/bin/
	chmod u+x ~/.local/bin/dlc
