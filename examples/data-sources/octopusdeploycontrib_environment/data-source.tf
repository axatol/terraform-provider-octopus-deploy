data "octopusdeploycontrib_environment" "by_name" {
  name = "Development"
}

data "octopusdeploycontrib_environment" "by_id" {
  id = "Environments-781"
}

data "octopusdeploycontrib_environment" "by_space_and_name" {
  space_id = "Spaces-142"
  name     = "SpinUp"
}
