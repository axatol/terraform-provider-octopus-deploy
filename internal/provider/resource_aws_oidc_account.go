package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/axatol/terraform-provider-octopusdeploycontrib/internal/custom"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                     = (*AWSOIDCAccountResource)(nil)
	_ resource.ResourceWithConfigValidators = (*AWSOIDCAccountResource)(nil)
	_ resource.ResourceWithConfigure        = (*AWSOIDCAccountResource)(nil)
	_ resource.ResourceWithImportState      = (*AWSOIDCAccountResource)(nil)
)

func NewAWSOIDCAccountResource() resource.Resource {
	return &AWSOIDCAccountResource{}
}

// AWSOIDCAccountResource defines the resource implementation.
type AWSOIDCAccountResource struct {
	client *client.Client
}

// AWSOIDCAccountResourceModel describes the resource data model.
type AWSOIDCAccountResourceModel struct {
	SpaceID                         types.String `tfsdk:"space_id"`
	ID                              types.String `tfsdk:"id"`
	Slug                            types.String `tfsdk:"slug"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	RoleARN                         types.String `tfsdk:"role_arn"`
	SessionDuration                 types.String `tfsdk:"session_duration"`
	EnvironmentIDs                  types.List   `tfsdk:"environment_ids"`
	TenantIDs                       types.List   `tfsdk:"tenant_ids"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`
	DeploymentSubjectKeys           types.List   `tfsdk:"deployment_subject_keys"`
	HealthCheckSubjectKeys          types.List   `tfsdk:"health_check_subject_keys"`
	AccountTestSubjectKeys          types.List   `tfsdk:"account_test_subject_keys"`
}

// expandAWSOIDCAccountResourceModel converts the model to a resource.
func expandAWSOIDCAccountResourceModel(ctx context.Context, model AWSOIDCAccountResourceModel) (*custom.AWSOIDCAccount, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := custom.AWSOIDCAccount{
		SpaceID:                         model.SpaceID.ValueString(),
		ID:                              model.ID.ValueString(),
		Slug:                            model.Slug.ValueString(),
		Name:                            model.Name.ValueString(),
		Description:                     model.Description.ValueString(),
		TenantedDeploymentParticipation: model.TenantedDeploymentParticipation.ValueString(),
		AccountType:                     "AmazonWebServicesOidcAccount",
		RoleARN:                         model.RoleARN.ValueString(),
		SessionDuration:                 model.SessionDuration.ValueString(),
	}

	var nestedDiags diag.Diagnostics
	resource.EnvironmentIDs, nestedDiags = expandStringList(ctx, model.EnvironmentIDs)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	resource.TenantIDs, nestedDiags = expandStringList(ctx, model.TenantIDs)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	resource.TenantTags, nestedDiags = expandStringList(ctx, model.TenantTags)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	resource.DeploymentSubjectKeys, nestedDiags = expandStringList(ctx, model.DeploymentSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	resource.HealthCheckSubjectKeys, nestedDiags = expandStringList(ctx, model.HealthCheckSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	resource.AccountTestSubjectKeys, nestedDiags = expandStringList(ctx, model.AccountTestSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &resource, diags
	}

	return &resource, diags
}

// flattenAWSOIDCAccountResourceModel converts the resource to a model.
func flattenAWSOIDCAccountResourceModel(ctx context.Context, resource *custom.AWSOIDCAccount) (*AWSOIDCAccountResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := AWSOIDCAccountResourceModel{
		SpaceID:                         types.StringValue(resource.SpaceID),
		ID:                              types.StringValue(resource.ID),
		Slug:                            types.StringValue(resource.Slug),
		Name:                            types.StringValue(resource.Name),
		Description:                     types.StringValue(resource.Description),
		TenantedDeploymentParticipation: types.StringValue(resource.TenantedDeploymentParticipation),
		RoleARN:                         types.StringValue(resource.RoleARN),
		SessionDuration:                 types.StringValue(resource.SessionDuration),
	}

	var nestedDiags diag.Diagnostics
	model.EnvironmentIDs, nestedDiags = flattenStringList(ctx, resource.EnvironmentIDs)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	model.TenantIDs, nestedDiags = flattenStringList(ctx, resource.TenantIDs)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	model.TenantTags, nestedDiags = flattenStringList(ctx, resource.TenantTags)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	model.DeploymentSubjectKeys, nestedDiags = flattenStringList(ctx, resource.DeploymentSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	model.HealthCheckSubjectKeys, nestedDiags = flattenStringList(ctx, resource.HealthCheckSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	model.AccountTestSubjectKeys, nestedDiags = flattenStringList(ctx, resource.AccountTestSubjectKeys)
	if diags.Append(nestedDiags...); diags.HasError() {
		return &model, diags
	}

	return &model, diags
}

func (r *AWSOIDCAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_aws_oidc_account"
}

func (r *AWSOIDCAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "The AWS OIDC account resource.",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				Description:   "The space ID.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"id": schema.StringAttribute{
				Description:   "The ID of the account.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"slug": schema.StringAttribute{
				Description:   "The slug of the account.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "The name of the account.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description:   "The description of the account.",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tenanted_deployment_participation": schema.StringAttribute{
				Description: "The tenanted deployment participation of the account.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf("Tenanted", "TenantedOrUntenanted", "Untenanted")},
			},
			"role_arn": schema.StringAttribute{
				Description: "The role ARN of the account.",
				Required:    true,
			},
			"session_duration": schema.StringAttribute{
				Description: "The session duration of the account.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("3600"),
			},
			"environment_ids": schema.ListAttribute{
				Description: "The environment IDs of the account.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"tenant_ids": schema.ListAttribute{
				Description: "The tenant IDs of the account.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"tenant_tags": schema.ListAttribute{
				Description: "The tenant tags of the account.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"deployment_subject_keys": schema.ListAttribute{
				Description: "Subject claims to include when using this account for deployments.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators:  []validator.List{listvalidator.ValueStringsAre(stringvalidator.OneOf("space", "environment", "project", "tenant", "runbook", "account", "type"))},
			},
			"health_check_subject_keys": schema.ListAttribute{
				Description: "Subject claims to include when using this account for health checks.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators:  []validator.List{listvalidator.ValueStringsAre(stringvalidator.OneOf("space", "account", "target", "type"))},
			},
			"account_test_subject_keys": schema.ListAttribute{
				Description: "Subject claims to include when using this account for account tests.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators:  []validator.List{listvalidator.ValueStringsAre(stringvalidator.OneOf("space", "account", "type"))},
			},
		},
	}
}

func (r *AWSOIDCAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *AWSOIDCAccountResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *AWSOIDCAccountResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan AWSOIDCAccountResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	resource, diags := expandAWSOIDCAccountResourceModel(ctx, plan)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if resource.SpaceID == "" {
		resource.SpaceID = r.client.GetSpaceID()
	}

	tflog.Debug(ctx, "creating resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource), "plan": fmt.Sprintf("%#v", plan)})

	resource, err := custom.NewClient(r.client).CreateAWSOIDCAccount(ctx, *resource)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to create resource", err)...); res.Diagnostics.HasError() {
		return
	}

	model, diags := flattenAWSOIDCAccountResourceModel(ctx, resource)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "created resource", map[string]interface{}{"resource": fmt.Sprintf("%#v", resource), "model": fmt.Sprintf("%#v", model)})

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *AWSOIDCAccountResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state AWSOIDCAccountResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	resourceID := state.ID.ValueString()
	spaceID := state.SpaceID.ValueString()

	tflog.Debug(ctx, "fetching resource", map[string]interface{}{"id": resourceID})

	resource, err := custom.NewClient(r.client).GetAWSOIDCAccount(ctx, spaceID, resourceID)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get resource", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched resource", map[string]interface{}{"resource": resource})

	model, diags := flattenAWSOIDCAccountResourceModel(ctx, resource)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *AWSOIDCAccountResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan AWSOIDCAccountResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	resource, diags := expandAWSOIDCAccountResourceModel(ctx, plan)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updating resource", map[string]interface{}{"resource": resource})

	resource, err := custom.NewClient(r.client).UpdateAWSOIDCAccount(ctx, *resource)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to update resource", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated resource", map[string]interface{}{"resource": resource})

	model, diags := flattenAWSOIDCAccountResourceModel(ctx, resource)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *AWSOIDCAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state AWSOIDCAccountResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	resource, diags := expandAWSOIDCAccountResourceModel(ctx, state)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleting resource", map[string]interface{}{"resource": resource})

	err := custom.NewClient(r.client).DeleteAWSOIDCAccount(ctx, resource.SpaceID, resource.ID)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to delete resource", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleted resource", map[string]interface{}{"resource": resource})
}

func (r *AWSOIDCAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	tflog.Debug(ctx, "importing resource", map[string]interface{}{"resource_id": req.ID})

	resource, err := custom.NewClient(r.client).GetAWSOIDCAccount(ctx, r.client.GetSpaceID(), req.ID)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get resource", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "imported resource", map[string]interface{}{"resource": resource})

	model, diags := flattenAWSOIDCAccountResourceModel(ctx, resource)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}
