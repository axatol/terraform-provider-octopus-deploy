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
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
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
		Attributes: map[string]schema.Attribute{
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

	var (
		environment *environments.Environment
		err         error
	)

	id := data.ID.ValueString()
	name := data.Name.ValueString()

	tflog.Debug(ctx, "fetching environment", map[string]interface{}{"id": id, "name": name})

	if id != "" {
		environment, err = d.client.Environments.GetByID(id)
	} else if name != "" {
		var environments []*environments.Environment
		environments, err = d.client.Environments.GetByName(name)
		for _, env := range environments {
			if env.Name == name {
				environment = env
			}
		}
	} else {
		err = fmt.Errorf("did not provide a valid identifier")
	}

	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch environment %s", id), err.Error())
		return
	}

	tflog.Debug(ctx, "fetched environment", map[string]interface{}{"environment": environment})

	model := EnvironmentDataSourceModel{
		ID:   types.StringValue(environment.ID),
		Name: types.StringValue(environment.Name),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
