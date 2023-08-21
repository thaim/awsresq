GIT_VERSION := $(shell git describe --abbrev=0 --tags)

ifndef GIT_VERSION
GIT_VERSION = main
endif

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binaries
	go build -ldflags "-s -w -X main.version=${GIT_VERSION}" -o bin/awsresq

test: ## Run test
	go test -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: mock
mock:
	go generate ./...
