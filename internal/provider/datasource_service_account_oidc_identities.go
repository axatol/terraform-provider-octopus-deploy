package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/axatol/terraform-provider-octopusdeploycontrib/internal/custom"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = (*ServiceAccountOIDCIdentities)(nil)
	_ datasource.DataSourceWithConfigure = (*ServiceAccountOIDCIdentities)(nil)
)

func NewServiceAccountOIDCIdentities() datasource.DataSource {
	return &ServiceAccountOIDCIdentities{}
}

// ServiceAccountOIDCIdentities defines the data source implementation.
type ServiceAccountOIDCIdentities struct {
	client *client.Client
}

// ServiceAccountOIDCIdentityModel describes the data source nested oidc identity data model.
type ServiceAccountOIDCIdentityModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Issuer  types.String `tfsdk:"issuer"`
	Subject types.String `tfsdk:"subject"`
}

// ServiceAccountOIDCIdentitiesModel describes the data source data model.
type ServiceAccountOIDCIdentitiesModel struct {
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	ExternalID       types.String `tfsdk:"external_id"`
	OIDCIdentities   types.List   `tfsdk:"oidc_identities"`
	Skip             types.Int64  `tfsdk:"skip"`
	Take             types.Int64  `tfsdk:"take"`
}

func (d *ServiceAccountOIDCIdentities) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_service_account_oidc_identities"
}

// Configure adds the provider configured client to the data source.
func (d *ServiceAccountOIDCIdentities) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

func (d *ServiceAccountOIDCIdentities) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the audience used to authenticate with this service account and the possible matching OIDC identity subjects",
		Attributes: map[string]schema.Attribute{
			"service_account_id": schema.StringAttribute{
				MarkdownDescription: "ID of the service account",
				Required:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The OIDC audience to use when requesting an access token",
				Computed:            true,
			},
			"oidc_identities": schema.ListNestedAttribute{
				MarkdownDescription: "List of OIDC identities which can authenticate with the associated service account",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the OIDC identity",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "A unique name for the OIDC identity within this service account",
							Computed:            true,
						},
						"issuer": schema.StringAttribute{
							MarkdownDescription: "The URL where OIDC tokens will be issued from. This must match exactly the issuer provided in the OIDC token",
							Computed:            true,
						},
						"subject": schema.StringAttribute{
							MarkdownDescription: "The subject of the OIDC identity. This must match exactly the subject provided in the OIDC token",
							Computed:            true,
						},
					},
				},
			},
			"skip": schema.Int64Attribute{
				MarkdownDescription: "Number of items to skip",
				Required:            true,
			},
			"take": schema.Int64Attribute{
				MarkdownDescription: "Number of items to takd",
				Required:            true,
			},
		},
	}
}

func (d *ServiceAccountOIDCIdentities) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var data ServiceAccountOIDCIdentitiesModel
	if res.Diagnostics.Append(req.Config.Get(ctx, &data)...); res.Diagnostics.HasError() {
		return
	}

	id := data.ServiceAccountID.ValueString()
	skip := data.Skip.ValueInt64()
	take := data.Take.ValueInt64()

	tflog.Debug(ctx, "fetching service account oidc identities", map[string]interface{}{"id": id})

	if id == "" {
		err := fmt.Errorf("did not provide a valid identifier")
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch project %s", id), err.Error())
		return
	}

	identities, err := custom.NewClient(d.client).ListServiceAccountOIDCIdentites(ctx, id, int(skip), int(take))
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch service account oidc identities %s", id), err.Error())
		return
	}

	tflog.Debug(ctx, "fetched service account oidc identities", map[string]interface{}{"identities": identities})

	oidcIdentities := []ServiceAccountOIDCIdentityModel{}
	for _, id := range identities.OIDCIdentities {
		oidcIdentities = append(oidcIdentities, ServiceAccountOIDCIdentityModel{
			ID:      types.StringValue(*id.ID),
			Name:    types.StringValue(id.Name),
			Issuer:  types.StringValue(id.Issuer),
			Subject: types.StringValue(id.Subject),
		})
	}

	oidcIdentitySchema, ok := req.Config.Schema.GetAttributes()["oidc_identities"].(schema.ListNestedAttribute)
	if !ok {
		err := fmt.Errorf("found invalid schema type for oidc_identities")
		res.Diagnostics.AddError(fmt.Sprintf("Failed to fetch service account oidc identities %s", id), err.Error())
		return
	}

	oidcIdentityList, diags := types.ListValueFrom(ctx, oidcIdentitySchema.NestedObject.Type(), oidcIdentities)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	model := ServiceAccountOIDCIdentitiesModel{
		ServiceAccountID: types.StringValue(id),
		ExternalID:       types.StringValue(identities.ExternalID),
		OIDCIdentities:   oidcIdentityList,
		Skip:             types.Int64Value(skip),
		Take:             types.Int64Value(take),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
