package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
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
		Attributes: map[string]schema.Attribute{
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

	id := data.ID.ValueString()
	if id == "" {
		id = data.Name.ValueString()
	}

	if id == "" {
		err := fmt.Errorf("did not provide a valid identifier")
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch tenant %s", id), err.Error())
		return
	}

	tenant, err := d.client.Tenants.GetByIdentifier(id)
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch tenant %s", id), err.Error())
		return
	}

	model := TenantDataSourceModel{
		ID:   types.StringValue(tenant.ID),
		Name: types.StringValue(tenant.Name),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
