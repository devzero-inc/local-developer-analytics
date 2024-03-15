SHELL := /bin/sh
BIN := /usr/local/bin

REV := $(shell git rev-parse HEAD)
CHANGES := $(shell test -n "$$(git status --porcelain)" && echo '+CHANGES' || true)

PACKAGE = lda
TARGET = lda

VERSION?=0.2.0
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# figure out current platform
UNAME := $(shell uname | tr A-Z a-z )

# docker registry prefix:
DOCKER_REGISTRY := registry.gitlab.codilas.com/codilas/devzero/lda

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}${CHANGES} -X main.Branch=${BRANCH}

CGO_CFLAGS=-I/usr/local/include
CGO_LDFLAGS=-L/usr/local/lib

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# count processors
NPROCS:=1
OS:=$(shell uname -s)
ifeq ($(OS),Linux)
	NPROCS:=$(shell grep -c ^processor /proc/cpuinfo)
endif
ifeq ($(OS),Darwin) # Assume Mac OS X
	NPROCS:=$(shell sysctl -n hw.ncpu)
endif

.DEFAULT_GOAL := all

.PHONY: \
	all \
	build \
	deps \
	doc \
	init \
	run \
	clean \
	lint \
	hooks \
	hooks-install \
	docker-create-buildx \
	docker-push-buildx

all: build debug

## Build binary
build:
	CGO_ENABLED=1 GOOS=$(UNAME) go build -a -tags netgo -ldflags="$(LDFLAGS)" -o "$(TARGET)" .

## Build and debug
debug: build
	./$(TARGET)

## Install tools and deps
deps: tools
	go get -t -v ./...

## Run local godoc server on :8080
doc:
	godoc -http=localhost:8080 -index

## Install precommit hooks
hooks-install:
	pre-commit install

## Run precommit hooks
hooks:
	pre-commit run -a

## Install tools, deps, build, check
init: deps buf build

## Lint code
lint:
	go fmt ./...

## Build and run
run: build
	./$(TARGET)

## Run go tests
test:
	go test ./... -v

## Clean up all generated files
clean:
	rm -rf ./docs/*
	rm -rf ./$(TARGET)

## Rebuild dockerfile and run docker compose
docker-compose:
	docker-compose build
	docker-compose up

## Rebuild dockerfile and run docker compose - ubuntu
docker-compose-ubuntu:
	docker-compose -f ubuntu.docker-compose.yml build
	docker-compose -f ubuntu.docker-compose.yml up

# ## Push docker image to registry
docker-push:
	$(eval NAME=${DOCKER_REGISTRY}/${TARGET})
	$(eval DOCKER_FQDN=${NAME}:${BRANCH}-${COMMIT})
	docker build -f Dockerfile . -t ${DOCKER_FQDN}
	docker push ${DOCKER_FQDN}
	@echo "Pushed  ${DOCKER_FQDN} !"

## Create buildx remote builder for m1 support
docker-create-buildx:
	export DOCKER_HOST=tcp://10.10.150.66:2375
	docker buildx create --driver docker-container --platform linux/amd64,linux/arm64 --name remote-builder
	docker buildx use remote-builder

## Push docker image to registry with buildx
docker-push-buildx:
	$(eval NAME=${DOCKER_REGISTRY}/${TARGET})
	$(eval DOCKER_FQDN=${NAME}:${BRANCH}-${COMMIT})
	docker buildx build --platform linux/amd64,linux/arm64 -t ${DOCKER_FQDN} . --push
	@echo "Pushed  ${DOCKER_FQDN} !"

# Fancy help message
# Source: https://gist.github.com/prwhite/8168133
# COLOR
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
