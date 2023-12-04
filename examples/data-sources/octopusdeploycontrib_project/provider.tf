terraform {
  required_providers {
    octopusdeploycontrib = {
      source = "registry.terraform.io/axatol/octopusdeploycontrib"
    }
  }
}

provider "octopusdeploycontrib" {
  server_url = "https://samples.octopus.app"
  space_id   = "Spaces-105"
  api_key    = "API-GUEST"
}
