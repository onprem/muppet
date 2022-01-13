# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.5.2. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for oapi-codegen variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(OAPI_CODEGEN)
#	@echo "Running oapi-codegen"
#	@$(OAPI_CODEGEN) <flags/args..>
#
OAPI_CODEGEN := $(GOBIN)/oapi-codegen-v1.9.0
$(OAPI_CODEGEN): $(BINGO_DIR)/oapi-codegen.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/oapi-codegen-v1.9.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=oapi-codegen.mod -o=$(GOBIN)/oapi-codegen-v1.9.0 "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"

