package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = (*EnvironmentDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*EnvironmentDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*EnvironmentDataSource)(nil)
)

func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

// EnvironmentDataSource defines the data source implementation.
type EnvironmentDataSource struct {
	client *client.Client
}

// EnvironmentDataSourceModel describes the data source data model.
type EnvironmentDataSourceModel struct {
	SpaceID types.String `tfsdk:"space_id"`
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
}

func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_environment"
}

// Configure adds the provider configured client to the data source.
func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.Append(ErrUnexpectedDataSourceConfigureType(req.ProviderData))
		return
	}

	d.client = client
}

func (d *EnvironmentDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{datasourcevalidator.AtLeastOneOf(
		path.MatchRoot("id"),
		path.MatchRoot("name"),
	)}
}

func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the ID of an environment",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				MarkdownDescription: "ID of the space",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the environment",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the environment",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (d *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data EnvironmentDataSourceModel
	if res.Diagnostics.Append(req.Config.Get(ctx, &data)...); res.Diagnostics.HasError() {
		return
	}

	spaceID := data.SpaceID.ValueString()
	name := data.Name.ValueString()
	id := data.ID.ValueString()
	query := environments.EnvironmentsQuery{
		Name: name,
		IDs:  []string{id},
		Skip: 0,
		Take: 1,
	}

	identifier := id
	if name != "" {
		identifier = name
	}

	tflog.Debug(ctx, "fetched environment", map[string]interface{}{"environment_identifier": identifier, "space_id": spaceID})

	resources, err := environments.Get(d.client, spaceID, query)
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch environment %s", identifier), err.Error())
		return
	}

	if len(resources.Items) < 1 {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch environment %s", identifier), "environment not found")
		return
	}

	resource := resources.Items[0]

	tflog.Debug(ctx, "fetched environment", map[string]interface{}{"environment": resource})

	model := EnvironmentDataSourceModel{
		SpaceID: types.StringValue(resource.SpaceID),
		ID:      types.StringValue(resource.ID),
		Name:    types.StringValue(resource.Name),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
