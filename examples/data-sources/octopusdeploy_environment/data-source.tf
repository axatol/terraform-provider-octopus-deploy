data "octopusdeploy_environment" "petclinic" {
  environment_name = "Development"
}

data "octopusdeploy_environment" "octopus_deploy" {
  environment_id = "Environments-781"
}
