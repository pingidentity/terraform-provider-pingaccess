package application

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
	config "github.com/pingidentity/terraform-provider-pingaccess/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

// ApplicationResource is a helper function to simplify the provider implementation.
func ApplicationResource() resource.Resource {
	return &applicationResource{}
}

// applicationResource is the resource implementation.
type applicationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type applicationResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	AccessValidatorId                     types.Int64  `tfsdk:"access_validator_id"`
	AgentCacheInvalidatedExpiration       types.Int64  `tfsdk:"agent_cache_invalidated_expiration"`
	AgentCacheInvalidatedResponseDuration types.Int64  `tfsdk:"agent_cache_invalidated_response_duration"`
	AgentId                               types.Int64  `tfsdk:"agent_id"`
	AllowEmptyPathSegments                types.Bool   `tfsdk:"allow_empty_path_segments"`
	ApplicationType                       types.String `tfsdk:"application_type"`
	AuthenticationChallengePolicyId       types.String `tfsdk:"authentication_challenge_policy_id"`
	CaseSensitivePath                     types.Bool   `tfsdk:"case_sensitive_path"`
	ContextRoot                           types.String `tfsdk:"context_root"`
	DefaultAuthType                       types.String `tfsdk:"default_auth_type"`
	Description                           types.String `tfsdk:"description"`
	// Destination                           types.String `tfsdk:"destination"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	FallbackPostEncoding types.String `tfsdk:"fallback_post_encoding"`
	// IdentityMappingIds                    types.Map    `tfsdk:"identity_mapping_ids"`
	Issuer                types.String `tfsdk:"issuer"`
	LastModified          types.Int64  `tfsdk:"last_modified"`
	ManualOrderingEnabled types.Bool   `tfsdk:"manual_ordering_enabled"`
	Name                  types.String `tfsdk:"name"`
	// Policy                                types.Map    `tfsdk:"policy"`
	Realm             types.String `tfsdk:"realm"`
	RequireHTTPS      types.Bool   `tfsdk:"require_https"`
	ResourceOrder     types.Set    `tfsdk:"resource_order"`
	SidebandClientId  types.String `tfsdk:"sideband_client_id"`
	SiteId            types.Int64  `tfsdk:"site_id"`
	SpaSupportEnabled types.Bool   `tfsdk:"spa_support_enabled"`
	VirtualHostIds    types.Set    `tfsdk:"virtual_host_ids"`
	WebSessionId      types.Int64  `tfsdk:"web_session_id"`
	// IdentityMappingIds                    types.Map[string, int]              `tfsdk:"identity_mapping_ids"`
	// Policy                                types.Map[string, List[PolicyItem]] `tfsdk:"policy"`
}

// GetSchema defines the schema for the resource.
func (r *applicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	applicationResourceSchema(ctx, req, resp, false)
}

func applicationResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_validator_id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"agent_cache_invalidated_expiration": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"agent_cache_invalidated_response_duration": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"agent_id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"allow_empty_path_segments": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"application_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_challenge_policy_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"case_sensitive_path": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"context_root": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_auth_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "destination": schema.StringAttribute{
			// 	Computed: true,
			// 	Optional: true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fallback_post_encoding": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "identity_mapping_ids": schema.MapAttribute{
			// 	Computed: true,
			// 	Optional: true,
			// 	PlanModifiers: []planmodifier.Map{
			// 		mapplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"issuer": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"manual_ordering_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "policy": schema.MapAttribute{
			// 	Computed: true,
			// 	Optional: true,
			// 	PlanModifiers: []planmodifier.Map{
			// 		mapplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"realm": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_order": schema.SetAttribute{
				ElementType: types.Int64Type,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"require_https": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sideband_client_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"spa_support_enabled": schema.BoolAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_host_ids": schema.SetAttribute{
				ElementType: types.Int64Type,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"web_session_id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"context_root", "name", "site_id", "virtual_hosts"})
	}
	resp.Schema = schema
}

func addOptionalApplicationFields(ctx context.Context, addRequest *client.Application, plan applicationResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	return nil
}

// func (r *applicationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
// 	var model applicationResource
// 	req.Plan.Get(ctx, &model)

// 	client.NewDefaultAuthTypeFromValue()

// 	req.Plan.GetAttribute(ctx,"default_auth_type",)
// }

// Metadata returns the resource type name.
func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readApplicationResponse(ctx context.Context, r *client.Application, state *applicationResourceModel, expectedValues *applicationResourceModel) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.AccessValidatorId = internaltypes.Int64TypeOrNil(r.AccessValidatorId)
	state.AgentCacheInvalidatedExpiration = internaltypes.Int64TypeOrNil(r.AgentCacheInvalidatedExpiration)
	state.AgentCacheInvalidatedResponseDuration = internaltypes.Int64TypeOrNil(r.AgentCacheInvalidatedResponseDuration)
	state.AgentId = internaltypes.Int64TypeOrNil(&r.AgentId)
	state.AllowEmptyPathSegments = internaltypes.BoolTypeOrNil(r.AllowEmptyPathSegments)
	r.GetApplicationType()
	// state.ApplicationType = applicationType
	state.AuthenticationChallengePolicyId = internaltypes.StringValueOrNull(r.AuthenticationChallengePolicyId)
	state.CaseSensitivePath = internaltypes.BoolTypeOrNil(r.CaseSensitivePath)
	state.ContextRoot = internaltypes.StringValueOrNull(r.ContextRoot)
	// state.DefaultAuthType = client.DefaultAuthType
	state.Description = internaltypes.StringValueOrNull(&r.Description)
	// state.Destination = types.DestinationView(r.Destination)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.FallbackPostEncoding = internaltypes.StringValueOrNull(r.FallbackPostEncoding)
	// state.IdentityMappingIds = types.Map[string, int](r.IdentityMappingIds)
	state.Issuer = internaltypes.StringValueOrNull(r.Issuer)
	state.LastModified = internaltypes.Int64TypeOrNil(r.LastModified)
	state.ManualOrderingEnabled = internaltypes.BoolTypeOrNil(r.ManualOrderingEnabled)
	state.Name = types.StringValue(r.Name)
	// state.Policy = types.Map[string, List[PolicyItem]](r.Policy)
	state.Realm = internaltypes.StringValueOrNull(r.Realm)
	state.RequireHTTPS = internaltypes.BoolTypeOrNil(r.RequireHTTPS)
	state.ResourceOrder = internaltypes.GetInt64Set(r.ResourceOrder)
	state.SidebandClientId = internaltypes.StringValueOrNull(r.SidebandClientId)
	state.SiteId = internaltypes.Int64TypeOrNil(&r.SiteId)
	state.SpaSupportEnabled = internaltypes.BoolTypeOrNil(&r.SpaSupportEnabled)
	state.VirtualHostIds = internaltypes.GetInt64Set(r.VirtualHostIds)
	state.WebSessionId = internaltypes.Int64TypeOrNil(r.WebSessionId)
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var virtualHosts []int64
	plan.VirtualHostIds.ElementsAs(ctx, &virtualHosts, false)
	createApplication := client.NewApplication(plan.Name.ValueString(),
		client.DefaultAuthType(plan.DefaultAuthType.ValueString()), plan.SpaSupportEnabled.ValueBool(),
		plan.ContextRoot.ValueString(), plan.SiteId.ValueInt64(), plan.AgentId.ValueInt64(), internaltypes.BaseTypesStringValueOrNull(plan.SidebandClientId),
		virtualHosts, internaltypes.BaseTypesStringValueOrNull(plan.AuthenticationChallengePolicyId))
	err := addOptionalApplicationFields(ctx, createApplication, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Application", err.Error())
		return
	}
	requestJson, err := createApplication.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateApplication := r.apiClient.ApplicationsApi.AddApplication(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateApplication = apiCreateApplication.Application(*createApplication)
	applicationResponse, httpResp, err := r.apiClient.ApplicationsApi.AddApplicationExecute(apiCreateApplication)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Application", err, httpResp)
		return
	}
	responseJson, err := applicationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state applicationResourceModel

	readApplicationResponse(ctx, applicationResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readApplication(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readApplication(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state applicationResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadApplication, httpResp, err := apiClient.ApplicationsApi.GetApplication(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a Application", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadApplication.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readApplicationResponse(ctx, apiReadApplication, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateApplication(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateApplication(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan applicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state applicationResourceModel
	req.State.Get(ctx, &state)
	UpdateApplication := apiClient.ApplicationsApi.UpdateApplication(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	var virtualHosts []int64
	plan.VirtualHostIds.ElementsAs(ctx, &virtualHosts, false)
	CreateUpdateRequest := client.NewApplication(plan.Name.ValueString(),
		client.DefaultAuthType(plan.DefaultAuthType.ValueString()), plan.SpaSupportEnabled.ValueBool(),
		plan.ContextRoot.ValueString(), plan.SiteId.ValueInt64(), plan.AgentId.ValueInt64(), plan.SidebandClientId.ValueString(),
		virtualHosts, plan.AuthenticationChallengePolicyId.ValueString())
	err := addOptionalApplicationFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Application", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateApplication = UpdateApplication.Application(*CreateUpdateRequest)
	updateApplicationResponse, httpResp, err := apiClient.ApplicationsApi.UpdateApplicationExecute(UpdateApplication)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Application", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateApplicationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readApplicationResponse(ctx, updateApplicationResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteApplication(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteApplication(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state applicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.ApplicationsApi.DeleteApplication(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Application", err, httpResp)
		return
	}

}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
