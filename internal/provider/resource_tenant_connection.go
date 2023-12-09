package provider

import (
	"context"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = (*TenantConnectionResource)(nil)
	_ resource.ResourceWithConfigure   = (*TenantConnectionResource)(nil)
	_ resource.ResourceWithImportState = (*TenantConnectionResource)(nil)
)

func NewTenantConnectionResource() resource.Resource {
	return &TenantConnectionResource{}
}

// TenantConnectionResource defines the resource implementation.
type TenantConnectionResource struct {
	client *client.Client
}

// TenantConnectionResourceModel describes the resource data model.
type TenantConnectionResourceModel struct {
	TenantID       types.String `tfsdk:"tenant_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
}

func (r *TenantConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_tenant_connection"
}

func (r *TenantConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this resource to connect a project to a tenant and environments",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant to connect to",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project to connect",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"environment_ids": schema.ListAttribute{
				MarkdownDescription: "list of applicable environments to connect",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *TenantConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *TenantConnectionResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan TenantConnectionResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	tenantID := plan.TenantID.ValueString()
	projectID := plan.ProjectID.ValueString()
	environmentIDVals := make([]types.String, 0, len(plan.EnvironmentIDs.Elements()))
	if res.Diagnostics.Append(plan.EnvironmentIDs.ElementsAs(ctx, &environmentIDVals, false)...); res.Diagnostics.HasError() {
		return
	}

	environmentIDs := make([]string, len(environmentIDVals))
	for i, val := range environmentIDVals {
		environmentIDs[i] = val.ValueString()
	}

	tflog.Debug(ctx, "fetching tenant", map[string]interface{}{"id": tenantID})

	tenant, err := r.client.Tenants.GetByID(tenantID)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched tenant, updating project environments", map[string]interface{}{"tenant": tenant})

	tenant.ProjectEnvironments[projectID] = environmentIDs
	_, err = r.client.Tenants.Update(tenant)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to update tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated tenant project environment", map[string]interface{}{"tenant": tenant})

	if res.Diagnostics.Append(res.State.Set(ctx, plan)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *TenantConnectionResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state TenantConnectionResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	tenantID := state.TenantID.ValueString()
	projectID := state.ProjectID.ValueString()

	tflog.Debug(ctx, "fetching tenant", map[string]interface{}{"id": tenantID})

	tenant, err := r.client.Tenants.GetByIdentifier(tenantID)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched tenant", map[string]interface{}{"tenant": tenant})

	environmentIDs, ok := tenant.ProjectEnvironments[projectID]
	if !ok {
		res.State.RemoveResource(ctx)
		return
	}

	slices.Sort(environmentIDs)
	environmentIDList, diags := types.ListValueFrom(ctx, types.StringType, environmentIDs)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	state = TenantConnectionResourceModel{
		TenantID:       types.StringValue(tenant.ID),
		ProjectID:      types.StringValue(projectID),
		EnvironmentIDs: environmentIDList,
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *TenantConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan TenantConnectionResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	tenantID := plan.TenantID.ValueString()
	projectID := plan.ProjectID.ValueString()
	environmentIDVals := make([]types.String, 0, len(plan.EnvironmentIDs.Elements()))
	if res.Diagnostics.Append(plan.EnvironmentIDs.ElementsAs(ctx, &environmentIDVals, false)...); res.Diagnostics.HasError() {
		return
	}

	environmentIDs := make([]string, len(environmentIDVals))
	for i, val := range environmentIDVals {
		environmentIDs[i] = val.ValueString()
	}

	tflog.Debug(ctx, "fetching tenant", map[string]interface{}{"id": tenantID})

	tenant, err := r.client.Tenants.GetByID(tenantID)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched tenant, updating project environments", map[string]interface{}{"tenant": tenant})

	tenant.ProjectEnvironments[projectID] = environmentIDs
	_, err = r.client.Tenants.Update(tenant)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to update tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated tenant project environments", map[string]interface{}{"tenant": tenant})

	if res.Diagnostics.Append(res.State.Set(ctx, plan)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *TenantConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state TenantConnectionResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	tenantID := state.TenantID.ValueString()
	projectID := state.ProjectID.ValueString()

	tflog.Debug(ctx, "fetching tenant", map[string]interface{}{"id": tenantID})

	tenant, err := r.client.Tenants.GetByIdentifier(tenantID)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched tenant, updating project environments", map[string]interface{}{"tenant": tenant})

	delete(tenant.ProjectEnvironments, projectID)
	_, err = r.client.Tenants.Update(tenant)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to update tenant", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated tenant project environments", map[string]interface{}{"tenant": tenant})
}

func (r *TenantConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 && len(parts) != 3 {
		res.Diagnostics.AddError(
			"Error importing tenant connection",
			"ID should be in the form tenant_id:project_id[:environment_id_1[+environment_id_2[+environment_id_n]]]",
		)
		return
	}

	tenantID := parts[0]
	projectID := parts[1]
	environmentIDs := []attr.Value{}
	if len(parts) == 3 {
		ids := strings.Split(parts[2], "+")
		for _, id := range ids {
			environmentIDs = append(environmentIDs, types.StringValue(id))
		}
	}

	environmentIDList, diags := types.ListValue(types.StringType, environmentIDs)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "imported tenant connection", map[string]interface{}{
		"tenant_id":       tenantID,
		"project_id":      projectID,
		"environment_ids": environmentIDs,
	})

	model := TenantConnectionResourceModel{
		TenantID:       types.StringValue(tenantID),
		ProjectID:      types.StringValue(projectID),
		EnvironmentIDs: environmentIDList,
	}

	if res.Diagnostics.Append(res.State.Set(ctx, &model)...); res.Diagnostics.HasError() {
		return
	}
}
