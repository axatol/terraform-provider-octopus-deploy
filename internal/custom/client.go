package custom

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
)

func NewClient(client *client.Client) *Client {
	return &Client{client}
}

type Client struct{ client *client.Client }
