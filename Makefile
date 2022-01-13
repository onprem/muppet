include .bingo/Variables.mk

.PHONY: generate
generate: pkg/api/server.go pkg/api/client.go pkg/api/types.go

pkg/api/server.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate chi-server -package api pkg/api/spec.yaml | gofmt -s > $@

pkg/api/client.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate client -package api pkg/api/spec.yaml | gofmt -s > $@

pkg/api/types.go: $(OAPI_CODEGEN) pkg/api/spec.yaml
	$(OAPI_CODEGEN) -generate types -package api pkg/api/spec.yaml | gofmt -s > $@
