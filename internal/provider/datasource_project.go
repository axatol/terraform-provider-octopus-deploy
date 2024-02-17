package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = (*ProjectDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*ProjectDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*ProjectDataSource)(nil)
)

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	client *client.Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	SpaceID types.String `tfsdk:"space_id"`
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Slug    types.String `tfsdk:"slug"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_project"
}

// Configure adds the provider configured client to the data source.
func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

func (d *ProjectDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{datasourcevalidator.AtLeastOneOf(
		path.MatchRoot("id"),
		path.MatchRoot("name"),
	)}
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the ID of a project",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the space that the project belongs to",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the project",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the project in Octopus Deploy. This name must be unique",
				Computed:            true,
				Optional:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "A human-readable, unique identifier, used to identify a project",
				Computed:            true,
			},
		},
	}
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data ProjectDataSourceModel
	if res.Diagnostics.Append(req.Config.Get(ctx, &data)...); res.Diagnostics.HasError() {
		return
	}

	spaceID := data.SpaceID.ValueString()
	name := data.Name.ValueString()
	id := data.ID.ValueString()
	query := projects.ProjectsQuery{
		Name: name,
		IDs:  []string{id},
		Skip: 0,
		Take: 1,
	}

	identifier := id
	if name != "" {
		identifier = name
	}

	tflog.Debug(ctx, "fetched project", map[string]interface{}{"project_identifier": identifier, "space_id": spaceID})

	resources, err := projects.Get(d.client, spaceID, query)
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch project %s", identifier), err.Error())
		return
	}

	if len(resources.Items) < 1 {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch project %s", identifier), "project not found")
		return
	}

	resource := resources.Items[0]

	tflog.Debug(ctx, "fetched project", map[string]interface{}{"project": resource})

	model := ProjectDataSourceModel{
		SpaceID: types.StringValue(resource.SpaceID),
		ID:      types.StringValue(resource.ID),
		Name:    types.StringValue(resource.Name),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
