resource "octopusdeploycontrib_aws_oidc_account" "aws" {
  name                              = "aws"
  role_arn                          = "arn:aws:iam::000000000000:role/some-role"
  tenanted_deployment_participation = "TenantedOrUntenanted"
  environment_ids                   = ["Environments-366"]                                         # sandbox
  tenant_ids                        = ["Tenants-398", "Tenants-390", "Tenants-388", "Tenants-389"] # regions

  # Deployment: space:[space-slug]:project:[project-slug]:tenant:[tenant-slug]:environment:[environment-slug]:account:[account-slug]
  # Runbook: space:[space-slug]:project:[project-slug]:runbook:[runbook-slug]:tenant:[tenant-slug]:environment:[environment-slug]:account:[account-slug]
  deployment_subject_keys = ["space", "account", "environment", "project", "tenant", "runbook"]

  # space:[space-slug]:account:[account-slug]:type:health
  health_check_subject_keys = ["space", "account", "type"]

  # space:[space-slug]:account:[account-slug]:type:health
  account_test_subject_keys = ["space", "account", "type"]
}
