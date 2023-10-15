package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type headerRoundTripper struct {
	transport http.RoundTripper
	apiKey    string
}

func (rt *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Octopus-ApiKey", rt.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return rt.transport.RoundTrip(req)
}

func NewClient(serverURL, apiKey, spaceID string) (*Client, error) {
	uri, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server url: %s", err)
	}

	transport := headerRoundTripper{
		transport: http.DefaultTransport,
		apiKey:    apiKey,
	}

	client := Client{
		client:    &http.Client{Transport: &transport},
		ServerURL: uri,
		SpaceID:   spaceID,
	}

	return &client, nil
}

type Client struct {
	client    *http.Client
	ServerURL *url.URL
	SpaceID   string
}

func (c *Client) endpoint(path string, query map[string]string) *url.URL {
	uri := c.ServerURL.JoinPath("/api")
	if c.SpaceID != "" {
		uri = uri.JoinPath(c.SpaceID)
	}
	uri = uri.JoinPath(path)

	q := uri.Query()
	for key, value := range query {
		q.Add(key, value)
	}

	uri.RawQuery = q.Encode()

	return uri
}

func (c *Client) do(ctx context.Context, method string, endpoint string, query map[string]string, payload any) ([]byte, *http.Response, error) {
	var body io.Reader = http.NoBody
	if payload != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal body: %s", err)
		}

		body = bytes.NewBuffer(raw)
	}

	uri := c.endpoint(endpoint, query).String()
	ctx = tflog.SetField(ctx, "method", method)
	ctx = tflog.SetField(ctx, "uri", uri)
	tflog.Info(ctx, "Making API request")

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %s", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to %s %s: %s", method, endpoint, err)
	}

	defer res.Body.Close()
	rawResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read request body: %s", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, res, fmt.Errorf("failed to %s %s: received unexpected status code %d - %s", method, endpoint, res.StatusCode, res.Status)
	}

	tflog.Info(ctx, "data", map[string]any{"raw": string(rawResponse)})

	return rawResponse, res, nil
}
