terraform {
  required_providers {
    octopusdeploy = {
      source = "registry.terraform.io/axatol/octopusdeploy"
    }
  }
}

provider "octopusdeploy" {
  server_url = "https://octopus.k8s.axatol.xyz"
  space_id   = "Spaces-1"
}
