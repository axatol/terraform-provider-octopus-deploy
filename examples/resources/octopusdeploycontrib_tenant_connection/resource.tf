data "octopusdeploycontrib_tenant" "apse2" {
  name = "test"
}

data "octopusdeploycontrib_project" "test_project" {
  name = "Test Project"
}

data "octopusdeploycontrib_environment" "development" {
  name = "Development"
}

data "octopusdeploycontrib_environment" "production" {
  name = "Production"
}

resource "octopusdeploycontrib_tenant_connection" "test" {
  tenant_id  = data.octopusdeploycontrib_tenant.apse2.id
  project_id = data.octopusdeploycontrib_project.test_project.id
  environment_ids = [
    data.octopusdeploycontrib_environment.development.id,
    data.octopusdeploycontrib_environment.production.id,
  ]
}
