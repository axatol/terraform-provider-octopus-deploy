package provider

import (
	"context"
	"fmt"

	"github.com/axatol/terraform-provider-octopusdeploy/internal/api"
	"github.com/axatol/terraform-provider-octopusdeploy/internal/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	client *api.Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	ProjectID           types.String `tfsdk:"project_id"`
	ProjectName         types.String `tfsdk:"project_name"`
	ProjectSlug         types.String `tfsdk:"project_slug"`
	SpaceID             types.String `tfsdk:"space_id"`
	ProjectGroupID      types.String `tfsdk:"project_group_id"`
	LifecycleID         types.String `tfsdk:"lifecycle_id"`
	VariableSetID       types.String `tfsdk:"variable_set_id"`
	DeploymentProcessID types.String `tfsdk:"deployment_process_id"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_project"
}

// Configure adds the provider configured client to the data source.
func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		res.Diagnostics.Append(ErrUnexpectedDataConfigureType(req.ProviderData))
		return
	}

	d.client = client
}

func (d *ProjectDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{datasourcevalidator.AtLeastOneOf(
		path.MatchRoot("project_id"),
		path.MatchRoot("project_name"),
	)}
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id":            schema.StringAttribute{Computed: true, Optional: true},
			"project_name":          schema.StringAttribute{Computed: true, Optional: true},
			"project_slug":          schema.StringAttribute{Computed: true},
			"space_id":              schema.StringAttribute{Computed: true},
			"project_group_id":      schema.StringAttribute{Computed: true},
			"lifecycle_id":          schema.StringAttribute{Computed: true},
			"variable_set_id":       schema.StringAttribute{Computed: true},
			"deployment_process_id": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data ProjectDataSourceModel
	if res.Diagnostics.Append(req.Config.Get(ctx, &data)...); res.Diagnostics.HasError() {
		return
	}

	var (
		id      string
		project *octopusdeploy.Project
		err     error
	)

	if id = data.ProjectID.ValueString(); id != "" {
		project, err = d.client.GetProjectByID(ctx, id)
	} else if id := data.ProjectName.ValueString(); id != "" {
		project, err = d.client.GetProjectByName(ctx, id)
	} else {
		err = fmt.Errorf("did not provide a valid identifier")
	}

	if project == nil {
		err = fmt.Errorf("no matching project found")
	}

	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch project %s", id), err.Error())
		return
	}

	model := ProjectDataSourceModel{
		ProjectID:           types.StringValue(project.ID),
		ProjectName:         types.StringValue(project.Name),
		ProjectSlug:         types.StringValue(project.Slug),
		SpaceID:             types.StringValue(project.SpaceID),
		ProjectGroupID:      types.StringValue(project.ProjectGroupID),
		LifecycleID:         types.StringValue(project.LifecycleID),
		VariableSetID:       types.StringValue(project.VariableSetID),
		DeploymentProcessID: types.StringValue(project.DeploymentProcessID),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
