package provider

import (
	"context"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/axatol/terraform-provider-octopusdeploycontrib/internal/custom"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = (*ServiceAccountOIDCIdentityResource)(nil)
	_ resource.ResourceWithConfigure   = (*ServiceAccountOIDCIdentityResource)(nil)
	_ resource.ResourceWithImportState = (*ServiceAccountOIDCIdentityResource)(nil)
)

func NewServiceAccountOIDCIdentity() resource.Resource {
	return &ServiceAccountOIDCIdentityResource{}
}

// ServiceAccountOIDCIdentityResource defines the resource implementation.
type ServiceAccountOIDCIdentityResource struct {
	client *client.Client
}

// ServiceAccountOIDCIdentityResourceModel describes the resource data model.
type ServiceAccountOIDCIdentityResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Issuer           types.String `tfsdk:"issuer"`
	Subject          types.String `tfsdk:"subject"`
}

func (r *ServiceAccountOIDCIdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_service_account_oidc_identity"
}

func (r *ServiceAccountOIDCIdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this resource to create and manage OIDC subject claims on a service account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the service account OIDC identity",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_account_id": schema.StringAttribute{
				MarkdownDescription: "ID of the service account to associate this identity to",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the identity",
				Required:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "OIDC issuer url",
				Required:            true,
			},
			"subject": schema.StringAttribute{
				MarkdownDescription: "OIDC subject claims",
				Required:            true,
			},
		},
	}
}

func (r *ServiceAccountOIDCIdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.Append(ErrUnexpectedResourceConfigureType(req.ProviderData))
		return
	}

	r.client = client
}

func (r *ServiceAccountOIDCIdentityResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan ServiceAccountOIDCIdentityResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	identity := custom.OIDCIdentity{
		ServiceAccountID: plan.ServiceAccountID.ValueString(),
		Name:             plan.Name.ValueString(),
		Issuer:           plan.Issuer.ValueString(),
		Subject:          plan.Subject.ValueString(),
	}

	tflog.Debug(ctx, "creating service account oidc identity", map[string]interface{}{"identity": identity})

	create, err := custom.NewClient(r.client).CreateServiceAccountOIDCIdentity(ctx, identity)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to create service account oidc identity", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "created service account oidc identity", map[string]interface{}{"response": create})

	plan.ID = types.StringValue(create.ID)
	if res.Diagnostics.Append(res.State.Set(ctx, plan)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ServiceAccountOIDCIdentityResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state ServiceAccountOIDCIdentityResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	identityID := state.ID.ValueString()
	serviceAccountID := state.ServiceAccountID.ValueString()

	tflog.Debug(ctx, "fetching service account oidc identity", map[string]interface{}{
		"identity_id":        identityID,
		"service_account_id": serviceAccountID,
	})

	identity, err := custom.NewClient(r.client).GetServiceAccountOIDCIdentity(ctx, serviceAccountID, identityID)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get service account oidc identity", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched service account oidc identity", map[string]interface{}{"identity": identity})

	state = ServiceAccountOIDCIdentityResourceModel{
		ID:               types.StringValue(*identity.ID),
		ServiceAccountID: types.StringValue(identity.ServiceAccountID),
		Name:             types.StringValue(identity.Name),
		Issuer:           types.StringValue(identity.Issuer),
		Subject:          types.StringValue(identity.Subject),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ServiceAccountOIDCIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan ServiceAccountOIDCIdentityResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	identity := custom.OIDCIdentity{
		ID:               plan.ID.ValueStringPointer(),
		ServiceAccountID: plan.ServiceAccountID.ValueString(),
		Name:             plan.Name.ValueString(),
		Issuer:           plan.Issuer.ValueString(),
		Subject:          plan.Subject.ValueString(),
	}

	tflog.Debug(ctx, "updating service account oidc identity", map[string]interface{}{"identity": identity})

	_, err := custom.NewClient(r.client).UpdateServiceAccountOIDCIdentity(ctx, identity)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to create service account oidc identity", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated service account oidc identity", map[string]interface{}{"identity": identity})

	plan = ServiceAccountOIDCIdentityResourceModel{
		ID:               types.StringValue(plan.ID.ValueString()),
		ServiceAccountID: types.StringValue(plan.ServiceAccountID.ValueString()),
		Name:             types.StringValue(plan.Name.ValueString()),
		Issuer:           types.StringValue(plan.Issuer.ValueString()),
		Subject:          types.StringValue(plan.Subject.ValueString()),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, plan)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ServiceAccountOIDCIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state ServiceAccountOIDCIdentityResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	identityID := state.ID.ValueString()
	serviceAccountID := state.ServiceAccountID.ValueString()

	tflog.Debug(ctx, "deleting service account oidc identity", map[string]interface{}{
		"identity_id":        identityID,
		"service_account_id": serviceAccountID,
	})

	_, err := custom.NewClient(r.client).DeleteServiceAccountOIDCIdentity(ctx, serviceAccountID, identityID)
	if !isAPIErrorNotFound(err) {
		res.Diagnostics.Append(ErrAsDiagnostic("Failed to delete service account oidc identity", err)...)
		return
	}

	res.State.RemoveResource(ctx)
	tflog.Debug(ctx, "deleted service account oidc identity")
}

func (r *ServiceAccountOIDCIdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		res.Diagnostics.AddError(
			"Error importing service account oidc identity",
			"ID should be in the form user_id:service_account_oidc_identity_id",
		)
		return
	}

	identityID := parts[1]
	serviceAccountID := parts[0]

	tflog.Debug(ctx, "importing service account oidc identity", map[string]interface{}{
		"identity_id":        identityID,
		"service_account_id": serviceAccountID,
	})

	identity, err := custom.NewClient(r.client).GetServiceAccountOIDCIdentity(ctx, serviceAccountID, identityID)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get service account oidc identity", err)...); err != nil {
		return
	}

	tflog.Debug(ctx, "imported service account oidc identity", map[string]interface{}{"identity": identity})

	model := ServiceAccountOIDCIdentityResourceModel{
		ID:               types.StringValue(*identity.ID),
		ServiceAccountID: types.StringValue(identity.ServiceAccountID),
		Name:             types.StringValue(identity.Name),
		Issuer:           types.StringValue(identity.Issuer),
		Subject:          types.StringValue(identity.Subject),
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
