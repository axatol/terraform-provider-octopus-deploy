GENSPEC_URL ?= "https://samples.octopus.app"
GENSPEC_CURL_ARGS := --header 'X-Octopus-ApiKey: API-GUEST'
GENSPEC_CURL_ARGS += --header 'Accept: application/json'
GENSPEC_CURL_ARGS += --header 'Content-Type: application/json'
GENSPEC_CMD_COMMENT ?= "This code is generated, DO NOT EDIT"
GENSPEC_CMD_ARGS := -packagecomment=$(GENSPEC_CMD_COMMENT) -packagename=octopusdeploy
GENSPEC_OUTPUT_DIR += ./internal/octopusdeploy

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

genspec:
ifeq (, $(shell which gojsonstruct))
	@echo installing gojsonstruct
	go install github.com/twpayne/go-jsonstruct/cmd/gojsonstruct@latest
endif

	curl $(GENSPEC_CURL_ARGS) $(GENSPEC_URL)/api/Spaces-1/projects/Projects-1465 \
	| gojsonstruct $(GENSPEC_CMD_ARGS) -typename=Project \
	> $(GENSPEC_OUTPUT_DIR)/project.gen.go 
