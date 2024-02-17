data "octopusdeploycontrib_project" "by_name" {
  name = "PetClinic"
}

data "octopusdeploycontrib_project" "by_id" {
  id = "Projects-861"
}

data "octopusdeploycontrib_project" "by_space_and_name" {
  space_id = "Spaces-142"
  name     = "Instance Infrastructure"
}
