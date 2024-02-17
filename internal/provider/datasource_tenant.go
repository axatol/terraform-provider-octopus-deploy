package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = (*TenantDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*TenantDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*TenantDataSource)(nil)
)

func NewTenantDataSource() datasource.DataSource {
	return &TenantDataSource{}
}

// TenantDataSource defines the data source implementation.
type TenantDataSource struct {
	client *client.Client
}

// TenantDataSourceModel describes the data source data model.
type TenantDataSourceModel struct {
	SpaceID types.String `tfsdk:"space_id"`
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
}

func (d *TenantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_tenant"
}

// Configure adds the provider configured client to the data source.
func (d *TenantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

func (d *TenantDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{datasourcevalidator.AtLeastOneOf(
		path.MatchRoot("id"),
		path.MatchRoot("name"),
	)}
}

func (d *TenantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the ID of a tenant",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				MarkdownDescription: "ID of the space",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (d *TenantDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data TenantDataSourceModel
	if res.Diagnostics.Append(req.Config.Get(ctx, &data)...); res.Diagnostics.HasError() {
		return
	}

	spaceID := data.SpaceID.ValueString()
	name := data.Name.ValueString()
	id := data.ID.ValueString()
	query := tenants.TenantsQuery{
		Name: name,
		IDs:  []string{id},
		Skip: 0,
		Take: 1,
	}

	identifier := id
	if name != "" {
		identifier = name
	}

	tflog.Debug(ctx, "fetched tenant", map[string]interface{}{"tenant_identifier": identifier, "space_id": spaceID})

	tenants, err := tenants.Get(d.client, spaceID, query)
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch tenant %s", identifier), err.Error())
		return
	}

	if len(tenants.Items) < 1 {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch tenant %s", identifier), "tenant not found")
		return
	}

	tenant := tenants.Items[0]

	tflog.Debug(ctx, "fetched tenant", map[string]interface{}{"tenant": tenant})

	model := TenantDataSourceModel{
		SpaceID: types.StringValue(tenant.SpaceID),
		ID:      types.StringValue(tenant.ID),
		Name:    types.StringValue(tenant.Name),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
