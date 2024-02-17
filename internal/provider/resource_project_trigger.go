package provider

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                     = (*ProjectTriggerResource)(nil)
	_ resource.ResourceWithConfigValidators = (*ProjectTriggerResource)(nil)
	_ resource.ResourceWithConfigure        = (*ProjectTriggerResource)(nil)
	_ resource.ResourceWithImportState      = (*ProjectTriggerResource)(nil)
)

func NewProjectTriggerResource() resource.Resource {
	return &ProjectTriggerResource{}
}

// ProjectTriggerResource defines the resource implementation.
type ProjectTriggerResource struct {
	client *client.Client
}

// ProjectTriggerResourceModel describes the resource data model.
type ProjectTriggerResourceModel struct {
	SpaceID     types.String `tfsdk:"space_id"`
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`

	RunRunbookAction *ProjectTriggerRunRunbookActionResourceModel `tfsdk:"run_runbook_action"`
	// TODO AutoDeployAction
	// TODO DeployLatestReleaseAction
	// TODO DeployNewReleaseAction
	// TODO CreateReleaseAction

	CronExpressionSchedule *ProjectTriggerCronExpressionScheduleResourceModel `tfsdk:"cron_expression_schedule"`
	// TODO DailySchedule
	// TODO DayOfMonthSchedule
	// TODO DateOfMonthSchedule
}

// ProjectTriggerRunbookActionResourceModel describes the runbook action data model.
type ProjectTriggerRunRunbookActionResourceModel struct {
	RunbookID      types.String `tfsdk:"runbook_id"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
	TenantIDs      types.List   `tfsdk:"tenant_ids"`
	TenantTags     types.List   `tfsdk:"tenant_tags"`
}

// ProjectTriggerCronExpressionScheduleResourceModel describes the cron expression schedule data model.
type ProjectTriggerCronExpressionScheduleResourceModel struct {
	CronExpression types.String `tfsdk:"cron_expression"`
	Timezone       types.String `tfsdk:"timezone"`
}

// expandProjectTriggerResourceModel converts the model to a resource.
func expandProjectTriggerResourceModel(ctx context.Context, model ProjectTriggerResourceModel) (*triggers.ProjectTrigger, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := &triggers.ProjectTrigger{
		SpaceID:     model.SpaceID.ValueString(),
		Resource:    resources.Resource{ID: model.ID.ValueString()},
		ProjectID:   model.ProjectID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsDisabled:  model.IsDisabled.ValueBool(),
	}

	var nestedDiags diag.Diagnostics
	if model.RunRunbookAction != nil {
		action := actions.NewRunRunbookAction()
		action.Runbook = model.RunRunbookAction.RunbookID.ValueString()

		action.Environments, nestedDiags = expandStringList(ctx, model.RunRunbookAction.EnvironmentIDs)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

		action.Tenants, nestedDiags = expandStringList(ctx, model.RunRunbookAction.TenantIDs)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

		action.TenantTags, nestedDiags = expandStringList(ctx, model.RunRunbookAction.TenantTags)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

		resource.Action = action
	}

	if model.CronExpressionSchedule != nil {
		resource.Filter = filters.NewCronScheduledTriggerFilter(
			model.CronExpressionSchedule.CronExpression.ValueString(),
			model.CronExpressionSchedule.Timezone.ValueString(),
		)
	}

	return resource, diags
}

// flattenProjectTriggerResourceModel converts the resource to a model.
func flattenProjectTriggerResourceModel(ctx context.Context, resource *triggers.ProjectTrigger) (*ProjectTriggerResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := ProjectTriggerResourceModel{
		SpaceID:     types.StringValue(resource.SpaceID),
		ID:          types.StringValue(resource.ID),
		ProjectID:   types.StringValue(resource.ProjectID),
		Name:        types.StringValue(resource.Name),
		Description: types.StringValue(resource.Description),
		IsDisabled:  types.BoolValue(resource.IsDisabled),
	}

	var nestedDiags diag.Diagnostics
	switch action := resource.Action.(type) {
	case *actions.RunRunbookAction:
		model.RunRunbookAction = &ProjectTriggerRunRunbookActionResourceModel{
			RunbookID: types.StringValue(action.Runbook),
		}

		model.RunRunbookAction.EnvironmentIDs, nestedDiags = flattenStringList(ctx, action.Environments)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

		model.RunRunbookAction.TenantIDs, nestedDiags = flattenStringList(ctx, action.Tenants)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

		model.RunRunbookAction.TenantTags, nestedDiags = flattenStringList(ctx, action.TenantTags)
		if diags.Append(nestedDiags...); diags.HasError() {
			return nil, diags
		}

	default:
		err := fmt.Errorf("unhandled action type: %s", resource.Action.GetActionType())
		diags.Append(ErrAsDiagnostic("Unhandled action type", err)...)
		return nil, diags
	}

	switch filter := resource.Filter.(type) {
	case *filters.CronScheduledTriggerFilter:
		model.CronExpressionSchedule = &ProjectTriggerCronExpressionScheduleResourceModel{
			CronExpression: types.StringValue(filter.CronExpression),
			Timezone:       types.StringValue(filter.TimeZone),
		}

	default:
		err := fmt.Errorf("unhandled filter type: %s", resource.Filter.GetFilterType())
		diags.Append(ErrAsDiagnostic("Unhandled filter type", err)...)
		return nil, diags
	}

	return &model, diags
}

func (r *ProjectTriggerResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_project_trigger"
}

func (r *ProjectTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Use this resource to create and manage scheduled triggers for runbooks",
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the space that the trigger is associated with",
				Computed:            true,
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the trigger",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the project that the trigger is associated with",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the trigger",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the trigger",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the trigger is disabled",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"run_runbook_action": schema.SingleNestedAttribute{
				MarkdownDescription: "An action to execute a runbook",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"runbook_id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier of the runbook that the trigger is associated with",
						Required:            true,
					},
					"environment_ids": schema.ListAttribute{
						MarkdownDescription: "The unique identifiers of the environments that the trigger is associated with",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
					"tenant_ids": schema.ListAttribute{
						MarkdownDescription: "The unique identifiers of the tenants that the trigger is associated with",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
					"tenant_tags": schema.ListAttribute{
						MarkdownDescription: "The tags of the tenants that the trigger is associated with",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"cron_expression_schedule": schema.SingleNestedAttribute{
				MarkdownDescription: "The cron expression schedule of the trigger",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"cron_expression": schema.StringAttribute{
						MarkdownDescription: "The cron expression that the trigger is scheduled to run",
						Required:            true,
					},
					"timezone": schema.StringAttribute{
						MarkdownDescription: "The timezone that the trigger is scheduled to run in",
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *ProjectTriggerResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *ProjectTriggerResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("run_runbook_action"),
			// TODO path.MatchRoot("auto_deploy_action"),
			// TODO path.MatchRoot("deploy_latest_release_action"),
			// TODO path.MatchRoot("deploy_new_release_action"),
			// TODO path.MatchRoot("create_release_action"),
		),
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("cron_expression_schedule"),
			// TODO path.MatchRoot("daily_schedule"),
			// TODO path.MatchRoot("day_of_month_schedule"),
			// TODO path.MatchRoot("date_of_month_schedule"),
		),
	}
}

func (r *ProjectTriggerResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan ProjectTriggerResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	trigger, diags := expandProjectTriggerResourceModel(ctx, plan)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating trigger", map[string]interface{}{"trigger": fmt.Sprintf("%#v", trigger), "plan": fmt.Sprintf("%#v", plan)})

	trigger, err := r.client.ProjectTriggers.Add(trigger)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to create trigger", err)...); res.Diagnostics.HasError() {
		return
	}

	model, diags := flattenProjectTriggerResourceModel(ctx, trigger)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "created trigger", map[string]interface{}{"trigger": fmt.Sprintf("%#v", trigger), "model": fmt.Sprintf("%#v", model)})

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ProjectTriggerResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state ProjectTriggerResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	triggerID := state.ID.ValueString()

	tflog.Debug(ctx, "fetching trigger", map[string]interface{}{"id": triggerID})

	trigger, err := r.client.ProjectTriggers.GetByID(triggerID)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get trigger", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "fetched trigger", map[string]interface{}{"trigger": trigger})

	model, diags := flattenProjectTriggerResourceModel(ctx, trigger)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ProjectTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan ProjectTriggerResourceModel
	if res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); res.Diagnostics.HasError() {
		return
	}

	trigger, diags := expandProjectTriggerResourceModel(ctx, plan)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updating trigger", map[string]interface{}{"trigger": trigger})

	trigger, err := r.client.ProjectTriggers.Update(trigger)
	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to update trigger", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "updated trigger", map[string]interface{}{"trigger": trigger})

	model, diags := flattenProjectTriggerResourceModel(ctx, trigger)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}

func (r *ProjectTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state ProjectTriggerResourceModel
	if res.Diagnostics.Append(req.State.Get(ctx, &state)...); res.Diagnostics.HasError() {
		return
	}

	trigger, diags := expandProjectTriggerResourceModel(ctx, state)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleting trigger", map[string]interface{}{"trigger": trigger})

	err := r.client.ProjectTriggers.Delete(trigger)
	if isAPIErrorNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to delete trigger", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleted trigger", map[string]interface{}{"trigger": trigger})
}

func (r *ProjectTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	tflog.Debug(ctx, "importing trigger", map[string]interface{}{"trigger_id": req.ID})

	trigger, err := r.client.ProjectTriggers.GetByID(req.ID)
	if isAPIErrorNotFound(err) {
		res.Diagnostics.Append(ErrAsDiagnostic("Trigger not found", err)...)
		return
	}

	if res.Diagnostics.Append(ErrAsDiagnostic("Failed to get trigger", err)...); res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "imported trigger", map[string]interface{}{"trigger": trigger})

	model, diags := flattenProjectTriggerResourceModel(ctx, trigger)
	if res.Diagnostics.Append(diags...); res.Diagnostics.HasError() {
		return
	}

	if res.Diagnostics.Append(res.State.Set(ctx, model)...); res.Diagnostics.HasError() {
		return
	}
}
