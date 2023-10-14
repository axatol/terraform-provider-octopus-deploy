// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OctopusDeployProvider satisfies various provider interfaces.
var _ provider.Provider = &OctopusDeployProvider{}

// OctopusDeployProvider defines the provider implementation.
type OctopusDeployProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OctopusDeployProviderModel describes the provider data model.
type OctopusDeployProviderModel struct {
	ServerURL types.String `tfsdk:"server_url"`
	APIKey    types.String `tfsdk:"api_key"`
}

func (p *OctopusDeployProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "octopusdeploy"
	resp.Version = p.version
}

func (p *OctopusDeployProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Description: "The URL of the Octopus Deploy REST API",
				Required:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key to use with the Octopus Deploy REST API",
				Required:    true,
			},
		},
	}
}

func (p *OctopusDeployProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OctopusDeployProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServerURL.IsNull() || data.ServerURL.ValueString() == "" {
		resp.Diagnostics.AddError("Must provide a valid server URL", "Server URL was empty")
		return
	}

	if data.APIKey.IsNull() || data.APIKey.ValueString() == "" {
		resp.Diagnostics.AddError("Must provide a valid API key", "API key was empty")
		return
	}

	client := newDataClient(data.APIKey.ValueString())
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OctopusDeployProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *OctopusDeployProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OctopusDeployProvider{
			version: version,
		}
	}
}

func newDataClient(apiKey string) *http.Client {
	return &http.Client{Transport: &authRoundTripper{
		transport: http.DefaultTransport,
		apiKey:    apiKey,
	}}
}

type authRoundTripper struct {
	transport http.RoundTripper
	apiKey    string
}

func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Octopus-ApiKey", rt.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return rt.transport.RoundTrip(req)
}
