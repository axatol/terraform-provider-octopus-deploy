package custom

import (
	"net/http"

	odclient "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/dghubble/sling"
)

func NewClient(client *odclient.Client) *Client {
	return &Client{client}
}

type Client struct{ client *odclient.Client }

func (c *Client) do(client *sling.Sling, output any) error {
	failure := new(core.APIError)
	res, err := client.Receive(output, failure)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
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

	return failure
}
