include .bingo/Variables.mk

SERVICE_BIN_NAME ?= muppet-service
AGENT_BIN_NAME   ?= command-agent

VERSION := $(strip $(shell [ -d .git ] && git describe --always --tags --dirty))
BUILD_DATE := $(shell date -u +"%Y-%m-%d")
VCS_BRANCH := $(strip $(shell git rev-parse --abbrev-ref HEAD))
TMP_DIR := $(shell pwd)/tmp
SERVICE_DOCKER_REPO ?= onprem/muppet-service
AGENT_DOCKER_REPO ?= onprem/command-agent

$(SERVICE_BIN_NAME): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 go build -a -ldflags '-s -w' -o $(SERVICE_BIN_NAME) ./cmd/service

$(AGENT_BIN_NAME): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 go build -a -ldflags '-s -w' -o $(AGENT_BIN_NAME) ./cmd/agent

.PHONY: build
build: $(SERVICE_BIN_NAME) $(AGENT_BIN_NAME)

.PHONY: build-service
build-service: $(SERVICE_BIN_NAME)

.PHONY: build-agent
build-agent: $(AGENT_BIN_NAME)

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test -v -race -short ./...

.PHONY: container
container: container-service container-agent

.PHONY: container-service
container-service: Dockerfile.service
	@docker build -f Dockerfile.service \
		-t $(SERVICE_DOCKER_REPO):$(VCS_BRANCH)-$(BUILD_DATE)-$(VERSION) .
	@docker tag $(SERVICE_DOCKER_REPO):$(VCS_BRANCH)-$(BUILD_DATE)-$(VERSION) $(SERVICE_DOCKER_REPO):latest

.PHONY: container-agent
container-agent: Dockerfile.agent
	@docker build -f Dockerfile.agent \
		-t $(AGENT_DOCKER_REPO):$(VCS_BRANCH)-$(BUILD_DATE)-$(VERSION) .
	@docker tag $(AGENT_DOCKER_REPO):$(VCS_BRANCH)-$(BUILD_DATE)-$(VERSION) $(AGENT_DOCKER_REPO):latest

.PHONY: generate
generate: pkg/api/server.go pkg/api/client.go pkg/api/types.go README.md

pkg/api/server.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate chi-server -package api pkg/api/spec.yaml | gofmt -s > $@

pkg/api/client.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate client -package api pkg/api/spec.yaml | gofmt -s > $@

pkg/api/types.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate types -package api pkg/api/spec.yaml | gofmt -s > $@

$(TMP_DIR)/help-agent.txt: $(AGENT_BIN_NAME) $(TMP_DIR)
	./$(AGENT_BIN_NAME) --help &> $(TMP_DIR)/help-agent.txt || true

README.md: $(EMBEDMD) $(TMP_DIR)/help-agent.txt
	$(EMBEDMD) -w README.md


$(TMP_DIR):
	mkdir -p $(TMP_DIR)
