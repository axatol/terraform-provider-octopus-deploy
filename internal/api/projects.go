package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/axatol/terraform-provider-octopusdeploy/internal/octopusdeploy"
)

func (c *Client) GetProjectByID(ctx context.Context, id string) (*octopusdeploy.Project, error) {
	endpoint := fmt.Sprintf("/projects/%s", id)
	raw, _, err := c.do(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project with id %s: %s", id, err)
	}

	var project octopusdeploy.Project
	if err := json.Unmarshal(raw, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data for project with id %s: %s", id, err)
	}

	return &project, nil
}

func (c *Client) GetProjectByName(ctx context.Context, name string) (*octopusdeploy.Project, error) {
	raw, _, err := c.do(ctx, http.MethodGet, "/projects", map[string]string{"name": name, "take": "1"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project with name %s: %s", name, err)
	}

	var projectList octopusdeploy.List[octopusdeploy.Project]
	if err := json.Unmarshal(raw, &projectList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data for project with name %s: %s", name, err)
	}

	if len(projectList.Items) < 1 {
		return nil, fmt.Errorf("failed to fetch project with name %s: no matching project found", name)
	}

	return &projectList.Items[0], nil
}
