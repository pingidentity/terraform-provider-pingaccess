package highAvailabilityProfiles

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
	config "github.com/pingidentity/terraform-provider-pingaccess/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &availabilityProfileResource{}
	_ resource.ResourceWithConfigure   = &availabilityProfileResource{}
	_ resource.ResourceWithImportState = &availabilityProfileResource{}
)

// EngineListenerResource is a helper function to simplify the provider implementation.
func AvailabilityProfileResource() resource.Resource {
	return &availabilityProfileResource{}
}

// engineListenerResource is the resource implementation.
type availabilityProfileResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type availabilityProfileResourceModel struct {
	ClassName     types.String `tfsdk:"classname"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Configuration types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *availabilityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	availabilityProfileResourceSchema(ctx, req, resp, false)
}

func availabilityProfileResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages HighAvailabilityProfiles",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"classname": schema.StringAttribute{
				Required: true,
			},
			"configuration": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"connect_timeout": schema.Float64Attribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"pooled_connection_timeout": schema.Float64Attribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"read_timeout": schema.Float64Attribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"max_retries": schema.Float64Attribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"retry_delay": schema.Float64Attribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"failed_retry_timeout": schema.Float64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Float64{
							float64planmodifier.UseStateForUnknown(),
						},
					},
					"failure_http_status_codes": schema.SetAttribute{
						Computed:    true,
						Optional:    true,
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}

	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"classname", "name", "configuration"})
	}
	resp.Schema = schema
}

func addOptionalAvailabilityProfileFields(ctx context.Context, addRequest *client.AvailabilityProfile, plan availabilityProfileResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsNonEmptyObj(plan.Configuration) {
		addRequest.Configuration = internaltypes.ObjValuesToMapNoPointer(plan.Configuration)
	}
	return nil
}

// Metadata returns the resource type name.
func (r *availabilityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_high_availability_profile"
}

func (r *availabilityProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readAvailabilityProfileResponse(ctx context.Context, r *client.AvailabilityProfile, state *availabilityProfileResourceModel, expectedValues *availabilityProfileResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.Name = types.StringValue(r.Name)
	state.ClassName = types.StringValue(r.ClassName)
	attrTypes := map[string]attr.Type{
		"connect_timeout":           basetypes.Float64Type{},
		"pooled_connection_timeout": basetypes.Float64Type{},
		"read_timeout":              basetypes.Float64Type{},
		"max_retries":               basetypes.Float64Type{},
		"retry_delay":               basetypes.Float64Type{},
		"failed_retry_timeout":      basetypes.Float64Type{},
		"failure_http_status_codes": basetypes.SetType{ElemType: types.StringType},
	}

	configValues := r.GetConfiguration()
	attrValues := map[string]attr.Value{
		"connect_timeout":           internaltypes.InterfaceFloat64TypeOrNull(configValues["connectTimeout"]),
		"pooled_connection_timeout": internaltypes.InterfaceFloat64TypeOrNull(configValues["pooledConnectionTimeout"]),
		"read_timeout":              internaltypes.InterfaceFloat64TypeOrNull(configValues["readTimeout"]),
		"max_retries":               internaltypes.InterfaceFloat64TypeOrNull(configValues["maxRetries"]),
		"retry_delay":               internaltypes.InterfaceFloat64TypeOrNull(configValues["retryDelay"]),
		"failed_retry_timeout":      internaltypes.InterfaceFloat64TypeOrNull(configValues["failedRetryTimeout"]),
		"failure_http_status_codes": internaltypes.GetInterfaceStringSet(configValues["failureHttpStatusCodes"]),
	}
	configuration := internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
	state.Configuration = configuration
}

func (r *availabilityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan availabilityProfileResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createAvailabilityProfile := client.NewAvailabilityProfile(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalAvailabilityProfileFields(ctx, createAvailabilityProfile, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for High Availability Profile", err.Error())
		return
	}
	requestJson, err := createAvailabilityProfile.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiCreateAvailabilityProfile := r.apiClient.HighAvailabilityApi.AddAvailabilityProfile(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAvailabilityProfile = apiCreateAvailabilityProfile.AvailabilityProfile(*createAvailabilityProfile)
	highAvailabilityProfileResponse, httpResp, err := r.apiClient.HighAvailabilityApi.AddAvailabilityProfileExecute(apiCreateAvailabilityProfile)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating High Availability Profile", err, httpResp)
		return
	}
	responseJson, err := highAvailabilityProfileResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state availabilityProfileResourceModel

	readAvailabilityProfileResponse(ctx, highAvailabilityProfileResponse, &state, &plan, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *availabilityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAvailabilityProfile(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAvailabilityProfile(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state availabilityProfileResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAvailabilityProfile, httpResp, err := apiClient.HighAvailabilityApi.GetAvailabilityProfile(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an High Availability Profile", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAvailabilityProfile.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAvailabilityProfileResponse(ctx, apiReadAvailabilityProfile, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *availabilityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAvailabilityProfile(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAvailabilityProfile(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan availabilityProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state availabilityProfileResourceModel
	req.State.Get(ctx, &state)
	UpdateAvailabilityProfile := apiClient.HighAvailabilityApi.UpdateAvailabilityProfile(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewAvailabilityProfile(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalAvailabilityProfileFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for High Availability Profile", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateAvailabilityProfile = UpdateAvailabilityProfile.AvailabilityProfile(*CreateUpdateRequest)
	UpdateAvailabilityProfileResponse, httpResp, err := apiClient.HighAvailabilityApi.UpdateAvailabilityProfileExecute(UpdateAvailabilityProfile)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating High Availability Profile", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := UpdateAvailabilityProfileResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readAvailabilityProfileResponse(ctx, UpdateAvailabilityProfileResponse, &state, &plan, &resp.Diagnostics)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *availabilityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAvailabilityProfile(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAvailabilityProfile(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state availabilityProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.HighAvailabilityApi.DeleteAvailabilityProfile(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting High Availability Profile", err, httpResp)
		return
	}
}

func (r *availabilityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
