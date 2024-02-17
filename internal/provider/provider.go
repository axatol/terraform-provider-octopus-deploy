// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/url"
	"os"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure OctopusDeployProvider satisfies various provider interfaces.
var (
	_ provider.Provider = (*OctopusDeployProvider)(nil)
)

// OctopusDeployProvider defines the provider implementation.
type OctopusDeployProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OctopusDeployProviderModel describes the provider data model.
type OctopusDeployProviderModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	ServerURL types.String `tfsdk:"server_url"`
	APIKey    types.String `tfsdk:"api_key"`
}

func (p *OctopusDeployProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "octopusdeploycontrib"
	resp.Version = p.version
}

func (p *OctopusDeployProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				MarkdownDescription: "The default space ID. Can be set with the environment variable `OCTOPUSDEPLOY_SPACE_ID`",
				Optional:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Octopus Deploy REST API. Can be set with the environment variable `OCTOPUSDEPLOY_SERVER_URL`",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key to use with the Octopus Deploy REST API. Can be set with the environment variable `OCTOPUSDEPLOY_API_KEY`",
				Optional:            true,
				Sensitive:           true,
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

	if data.SpaceID.IsUnknown() {
		resp.Diagnostics.Append(ErrUnknownProviderAttribute("space_id", "OCTOPUSDEPLOY_SPACE_ID"))
	}

	if data.ServerURL.IsUnknown() {
		resp.Diagnostics.Append(ErrUnknownProviderAttribute("server_url", "OCTOPUSDEPLOY_SERVER_URL"))
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.Append(ErrUnknownProviderAttribute("api_key", "OCTOPUSDEPLOY_API_KEY"))
	}

	spaceID := os.Getenv("OCTOPUSDEPLOY_SPACE_ID")
	serverURL := os.Getenv("OCTOPUSDEPLOY_SERVER_URL")
	apiKey := os.Getenv("OCTOPUSDEPLOY_API_KEY")
	ctx = tflog.SetField(ctx, "space_id", spaceID)
	ctx = tflog.SetField(ctx, "server_url", serverURL)
	ctx = tflog.SetField(ctx, "api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_key")
	tflog.Info(ctx, "Retrieved provider data from environment variables")

	if !data.SpaceID.IsNull() {
		spaceID = data.SpaceID.ValueString()
	}

	if !data.ServerURL.IsNull() {
		serverURL = data.ServerURL.ValueString()
	}

	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	if spaceID == "" {
		resp.Diagnostics.Append(ErrMissingProviderAttribute("space_url", "OCTOPUSDEPLOY_SPACE_ID"))
	}

	if serverURL == "" {
		resp.Diagnostics.Append(ErrMissingProviderAttribute("server_url", "OCTOPUSDEPLOY_SERVER_URL"))
	}

	if apiKey == "" {
		resp.Diagnostics.Append(ErrMissingProviderAttribute("api_key", "OCTOPUSDEPLOY_API_KEY"))
	}

	ctx = tflog.SetField(ctx, "space_id", spaceID)
	ctx = tflog.SetField(ctx, "server_url", serverURL)
	ctx = tflog.SetField(ctx, "api_key", apiKey)
	tflog.Info(ctx, "Provider configuration resolved")

	if resp.Diagnostics.HasError() {
		return
	}

	uri, err := url.Parse(serverURL)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse Octopus Deploy server URL", err.Error())
		return
	}

	client, err := client.NewClient(nil, uri, apiKey, spaceID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Octopus Deploy API client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OctopusDeployProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectTriggerResource,
		NewServiceAccountOIDCIdentity,
		NewTenantConnectionResource,
	}
}

func (p *OctopusDeployProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEnvironmentDataSource,
		NewProjectDataSource,
		NewServiceAccountOIDCIdentities,
		NewTenantDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OctopusDeployProvider{
			version: version,
		}
	}
}
