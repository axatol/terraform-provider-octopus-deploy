data "octopusdeploycontrib_tenant" "by_name" {
  name = "Brisbane Vet"
}

data "octopusdeploycontrib_tenant" "by_id" {
  id = "Tenants-381"
}

data "octopusdeploycontrib_tenant" "by_space_and_name" {
  space_id = "Spaces-142"
  name     = "Internal"
}
