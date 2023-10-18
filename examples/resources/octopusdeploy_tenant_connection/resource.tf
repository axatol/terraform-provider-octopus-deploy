data "octopusdeploy_tenant" "apse2" {
  name = "test"
}

data "octopusdeploy_project" "test_project" {
  name = "Test Project"
}

data "octopusdeploy_environment" "development" {
  name = "Development"
}

data "octopusdeploy_environment" "production" {
  name = "Production"
}

resource "octopusdeploy_tenant_connection" "test" {
  tenant_id  = data.octopusdeploy_tenant.apse2.id
  project_id = data.octopusdeploy_project.test_project.id
  environment_ids = [
    data.octopusdeploy_environment.development.id,
    data.octopusdeploy_environment.production.id,
  ]
}
