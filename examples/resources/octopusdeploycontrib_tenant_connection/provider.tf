terraform {
  required_providers {
    octopusdeploy = {
      source = "registry.terraform.io/OctopusDeployLabs/octopusdeploy"
    }

    octopusdeploycontrib = {
      source = "registry.terraform.io/axatol/octopusdeploycontrib"
    }
  }
}

provider "octopusdeploy" {
}

provider "octopusdeploycontrib" {
}
