resource "octopusdeploycontrib_project_trigger" "name" {
  name       = "test"
  project_id = "Projects-2"
  space_id   = "Spaces-1"

  cron_expression_schedule = {
    cron_expression = "1 1 * * 0"
    timezone        = "UTC"
  }

  run_runbook_action = {
    runbook_id      = "Runbooks-1"
    environment_ids = ["Environments-2"]
  }
}
