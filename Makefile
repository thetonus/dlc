SHELL := /bin/bash
BINARY_NAME := dlc
PACKAGE_NAMESPACE := github.com/hammacktony/dlc
GO := go
CGO := 0
GO111MODULE := on

all: help

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: mod
mod: ## Go mod things
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod download

build: clean ## Build app
	$(eval DLC_VERSION ?= $(shell git describe --tags --match 'v*' --abbrev=0)+dev)
	$(eval DLC_COMMIT ?= $(shell git rev-parse --short HEAD))

	CGO=$(CGO) GO111MODULE=$(GO111MODULE) $(GO) \
		build \
		-ldflags "-X $(PACKAGE_NAMESPACE)/pkg/global.Version=$(DLC_VERSION) -X $(PACKAGE_NAMESPACE)/pkg/global.Commit=$(DLC_COMMIT) -X $(PACKAGE_NAMESPACE)/pkg/global.BuildTime=$(shell date +%Y-%m-%dT%H:%M:%S%z) -w" \
		-o dist/${BINARY_NAME} cmd/main.go

clean: ## Clean stuff
	$(GO) clean
	rm -rf dist/*

format: ## Format
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

dev/install: build ## Install dlc locally
	rm -rf ~/.local/bin/dlc
	cp dist/dlc ~/.local/bin/
	chmod u+x ~/.local/bin/dlc
