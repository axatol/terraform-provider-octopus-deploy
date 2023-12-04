terraform {
  required_providers {
    octopusdeploycontrib = {
      source = "registry.terraform.io/axatol/octopusdeploycontrib"
    }
  }
}

provider "octopusdeploycontrib" {
  server_url = "https://octopus.k8s.axatol.xyz"
  space_id   = "Spaces-1"
}
