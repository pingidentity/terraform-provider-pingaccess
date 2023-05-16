package thirdPartyService

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
	config "github.com/pingidentity/terraform-provider-pingaccess/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &thirdPartyServiceResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyServiceResource{}
	_ resource.ResourceWithImportState = &thirdPartyServiceResource{}
)

// ThirdPartyServiceResource is a helper function to simplify the provider implementation.
func ThirdPartyServiceResource() resource.Resource {
	return &thirdPartyServiceResource{}
}

// thirdPartyServiceResource is the resource implementation.
type thirdPartyServiceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type thirdPartyServiceResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	AvailabilityProfileId     types.Int64  `tfsdk:"availability_profile_id"`
	ExpectedHostname          types.String `tfsdk:"expected_hostname"`
	HostValue                 types.String `tfsdk:"host_value"`
	LoadBalancingStrategyId   types.Int64  `tfsdk:"load_balancing_strategy_id"`
	MaxConnections            types.Int64  `tfsdk:"max_connections"`
	Name                      types.String `tfsdk:"name"`
	Secure                    types.Bool   `tfsdk:"secure"`
	SkipHostnameVerification  types.Bool   `tfsdk:"skip_hostname_verification"`
	Targets                   types.Set    `tfsdk:"targets"`
	TrustedCertificateGroupId types.Int64  `tfsdk:"trusted_certificate_group_id"`
	UseProxy                  types.Bool   `tfsdk:"use_proxy"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyServiceResourceSchema(ctx, req, resp, false)
}

func thirdPartyServiceResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a ThirdPartyService.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				// Add the other necessary attributes here
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"availability_profile_id": schema.Int64Attribute{
				Required: true,
			},
			"targets": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"expected_hostname": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host_value": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancing_strategy_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"secure": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_hostname_verification": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trusted_certificate_group_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"use_proxy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "targets"})
	}
	resp.Schema = schema
}
func addOptionalThirdPartyServiceFields(ctx context.Context, addRequest *client.ThirdPartyService, plan thirdPartyServiceResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AvailabilityProfileId) {
		intVal := plan.AvailabilityProfileId.ValueInt64()
		addRequest.AvailabilityProfileId = &intVal
	}
	if internaltypes.IsDefined(plan.ExpectedHostname) {
		stringVal := plan.ExpectedHostname.ValueString()
		addRequest.ExpectedHostname = &stringVal
	}
	if internaltypes.IsDefined(plan.HostValue) {
		stringVal := plan.HostValue.ValueString()
		addRequest.HostValue = &stringVal
	}
	if internaltypes.IsDefined(plan.LoadBalancingStrategyId) {
		intVal := plan.LoadBalancingStrategyId.ValueInt64()
		addRequest.LoadBalancingStrategyId = &intVal
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		intVal := plan.MaxConnections.ValueInt64()
		addRequest.MaxConnections = &intVal
	}
	if internaltypes.IsDefined(plan.TrustedCertificateGroupId) {
		intVal := plan.TrustedCertificateGroupId.ValueInt64()
		addRequest.TrustedCertificateGroupId = &intVal
	}
	if internaltypes.IsDefined(plan.Secure) {
		boolVal := plan.Secure.ValueBool()
		addRequest.Secure = &boolVal
	}
	if internaltypes.IsDefined(plan.SkipHostnameVerification) {
		boolVal := plan.SkipHostnameVerification.ValueBool()
		addRequest.SkipHostnameVerification = &boolVal
	}
	if internaltypes.IsDefined(plan.UseProxy) {
		boolVal := plan.UseProxy.ValueBool()
		addRequest.UseProxy = &boolVal
	}
	return nil
}

// Metadata returns the resource type name.
func (r *thirdPartyServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_services"
}

func (r *thirdPartyServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readThirdPartyServiceResponse(ctx context.Context, r *client.ThirdPartyService, state *thirdPartyServiceResourceModel, expectedValues *thirdPartyServiceResourceModel) {
	state.Id = types.StringValue(*r.Id)
	state.AvailabilityProfileId = types.Int64Value(*r.AvailabilityProfileId)
	state.ExpectedHostname = internaltypes.StringTypeOrNil(r.ExpectedHostname, false)
	state.HostValue = internaltypes.StringTypeOrNil(r.HostValue, false)
	state.LoadBalancingStrategyId = types.Int64Value(*r.LoadBalancingStrategyId)
	state.MaxConnections = types.Int64Value(*r.MaxConnections)
	state.Name = types.StringValue(r.Name)
	state.Secure = internaltypes.BoolTypeOrNil(r.Secure)
	state.SkipHostnameVerification = internaltypes.BoolTypeOrNil(r.SkipHostnameVerification)
	state.Targets = internaltypes.GetStringSet(r.Targets)
	state.TrustedCertificateGroupId = types.Int64Value(*r.TrustedCertificateGroupId)
	state.UseProxy = internaltypes.BoolTypeOrNil(r.UseProxy)
}

func (r *thirdPartyServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan thirdPartyServiceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var TargetsSlice []string
	plan.Targets.ElementsAs(ctx, &TargetsSlice, false)
	createThirdPartyService := client.NewThirdPartyService(TargetsSlice, plan.Name.ValueString())
	err := addOptionalThirdPartyServiceFields(ctx, createThirdPartyService, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ThirdPartyService", err.Error())
		return
	}
	requestJson, err := createThirdPartyService.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateThirdPartyService := r.apiClient.ThirdPartyServicesApi.AddThirdPartyService(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateThirdPartyService = apiCreateThirdPartyService.ThirdPartyService(*createThirdPartyService)
	thirdPartyServiceResponse, httpResp, err := r.apiClient.ThirdPartyServicesApi.AddThirdPartyServiceExecute(apiCreateThirdPartyService)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ThirdPartyService", err, httpResp)
		return
	}
	responseJson, err := thirdPartyServiceResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyServiceResourceModel

	readThirdPartyServiceResponse(ctx, thirdPartyServiceResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *thirdPartyServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyService(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyService(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state thirdPartyServiceResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadThirdPartyService, httpResp, err := apiClient.ThirdPartyServicesApi.GetThirdPartyService(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a ThirdPartyService", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadThirdPartyService.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyServiceResponse(ctx, apiReadThirdPartyService, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *thirdPartyServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyService(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyService(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyServiceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyServiceResourceModel
	req.State.Get(ctx, &state)
	var TargetsSlice []string
	plan.Targets.ElementsAs(ctx, &TargetsSlice, false)
	UpdateThirdPartyService := apiClient.ThirdPartyServicesApi.UpdateThirdPartyService(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewThirdPartyService(TargetsSlice, plan.Name.ValueString())
	err := addOptionalThirdPartyServiceFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ThirdPartyService", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateThirdPartyService = UpdateThirdPartyService.ThirdPartyService(*CreateUpdateRequest)
	updateThirdPartyServiceResponse, httpResp, err := apiClient.ThirdPartyServicesApi.UpdateThirdPartyServiceExecute(UpdateThirdPartyService)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating ThirdPartyService", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateThirdPartyServiceResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readThirdPartyServiceResponse(ctx, updateThirdPartyServiceResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *thirdPartyServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteThirdPartyService(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteThirdPartyService(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state thirdPartyServiceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.ThirdPartyServicesApi.DeleteThirdPartyService(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a ThirdPartyService", err, httpResp)
		return
	}

}

func (r *thirdPartyServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
