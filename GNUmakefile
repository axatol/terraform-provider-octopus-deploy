TF_LOG ?= INFO

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

install:
	go install .

generate:
	go generate

plan: install
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_environment plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_project plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_service_account_oidc_identities plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/data-sources/octopusdeploycontrib_tenant plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/resources/octopusdeploycontrib_aws_oidc_account plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/resources/octopusdeploycontrib_project_trigger plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/resources/octopusdeploycontrib_service_account_oidc_identity plan
	TF_LOG=$(TF_LOG) terraform -chdir=examples/resources/octopusdeploycontrib_tenant_connection plan
