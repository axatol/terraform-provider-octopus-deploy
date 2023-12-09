resource "octopusdeploycontrib_service_account_oidc_identity" "github" {
  service_account_id = "Users-41"
  name               = "GitHub"
  issuer             = "https://token.actions.githubusercontent.com"
  subject            = "repo:axatol/terraform-provider-octopusdeploycontrib:pull_request"
}
