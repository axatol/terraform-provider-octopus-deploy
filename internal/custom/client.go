package custom

import (
	"context"
	"fmt"
	"net/http"

	odclient "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewClient(client *odclient.Client) *Client {
	return &Client{client}
}

type Client struct{ client *odclient.Client }

func (c *Client) do(ctx context.Context, client *sling.Sling, output any) error {
	failure := new(core.APIError)
	res, err := client.Receive(output, failure)
	req, _ := client.Request()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		tflog.Debug(ctx, fmt.Sprintf("%s %s was successful", req.Method, req.URL.Path), map[string]interface{}{
			"status_code": res.StatusCode,
		})

		return nil
	}

	if failure.StatusCode == 0 {
		failure.StatusCode = res.StatusCode
	}

	if err != nil {
		failure.Errors = append(failure.Errors, err.Error())
	}

	if failure.ErrorMessage == "" {
		failure.ErrorMessage = http.StatusText(res.StatusCode)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s %s was unsuccessful", req.Method, req.URL.Path), map[string]interface{}{
		"status_code": res.StatusCode,
		"error":       fmt.Sprintf("%#v", failure),
	})

	return failure
}
