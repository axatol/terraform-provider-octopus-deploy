resource "octopusdeploy_tenant" "test" {
  name = "Test Tenant"
  lifecycle {
    ignore_changes = [project_environment]
  }
}

data "octopusdeploycontrib_project" "test" {
  name = "Test Project"
}

data "octopusdeploycontrib_environment" "development" {
  name = "Development"
}

data "octopusdeploycontrib_environment" "production" {
  name = "Production"
}

resource "octopusdeploycontrib_tenant_connection" "test" {
  tenant_id  = octopusdeploy_tenant.test.id
  project_id = data.octopusdeploycontrib_project.test.id
  environment_ids = [
    data.octopusdeploycontrib_environment.development.id,
    data.octopusdeploycontrib_environment.production.id,
  ]
}
