data "octopusdeploycontrib_service_account_oidc_identities" "github" {
  service_account_id = "Users-41"
  skip               = 0
  take               = 5
}

resource "terraform_data" "github" {
  input = data.octopusdeploycontrib_service_account_oidc_identities.github.oidc_identities
}
