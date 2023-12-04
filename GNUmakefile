TF_LOG ?= INFO

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

install:
	go install .

plan: install
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_environment plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_project plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_tenant plan
