package site

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
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// SiteResource is a helper function to simplify the provider implementation.
func SiteResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type siteResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	AvailabilityProfileId     types.Int64  `tfsdk:"availability_profile_id"`
	ExpectedHostname          types.String `tfsdk:"expected_hostname"`
	KeepAliveTimeout          types.Int64  `tfsdk:"keep_alive_timeout"`
	LoadBalancingStrategyId   types.Int64  `tfsdk:"load_balancing_strategy_id"`
	MaxConnections            types.Int64  `tfsdk:"max_connections"`
	MaxWebSocketConnections   types.Int64  `tfsdk:"max_web_socket_connections"`
	Name                      types.String `tfsdk:"name"`
	Secure                    types.Bool   `tfsdk:"secure"`
	SendPaCookie              types.Bool   `tfsdk:"send_pa_cookie"`
	SiteAuthenticatorIds      types.Set    `tfsdk:"site_authenticator_ids"`
	SkipHostnameVerification  types.Bool   `tfsdk:"skip_hostname_verification"`
	Targets                   types.Set    `tfsdk:"targets"`
	TrustedCertificateGroupId types.Int64  `tfsdk:"trusted_certificate_group_id"`
	UseProxy                  types.Bool   `tfsdk:"use_proxy"`
	UseTargetHostHeader       types.Bool   `tfsdk:"use_target_host_header"`
}

// GetSchema defines the schema for the resource.
func (r *siteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	siteResourceSchema(ctx, req, resp, false)
}

func siteResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"targets": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"availability_profile_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expected_hostname": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_alive_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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
			"max_web_socket_connections": schema.Int64Attribute{
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
			"send_pa_cookie": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"site_authenticator_ids": schema.SetAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.Int64Type,
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
			"use_target_host_header": schema.BoolAttribute{
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
func addOptionalSiteFields(ctx context.Context, addRequest *client.Site, plan siteResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsDefined(plan.SiteAuthenticatorIds) {
		var slice []int64
		plan.SiteAuthenticatorIds.ElementsAs(ctx, &slice, false)
		addRequest.SiteAuthenticatorIds = slice
	}
	if internaltypes.IsDefined(plan.ExpectedHostname) {
		stringVal := plan.ExpectedHostname.ValueString()
		addRequest.ExpectedHostname = &stringVal
	}
	if internaltypes.IsDefined(plan.AvailabilityProfileId) {
		intVal := plan.AvailabilityProfileId.ValueInt64()
		addRequest.AvailabilityProfileId = &intVal
	}
	if internaltypes.IsDefined(plan.KeepAliveTimeout) {
		intVal := plan.KeepAliveTimeout.ValueInt64()
		addRequest.KeepAliveTimeout = &intVal
	}
	if internaltypes.IsDefined(plan.LoadBalancingStrategyId) {
		intVal := plan.LoadBalancingStrategyId.ValueInt64()
		addRequest.LoadBalancingStrategyId = &intVal
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		intVal := plan.MaxConnections.ValueInt64()
		addRequest.MaxConnections = &intVal
	}
	if internaltypes.IsDefined(plan.MaxWebSocketConnections) {
		intVal := plan.MaxWebSocketConnections.ValueInt64()
		addRequest.MaxWebSocketConnections = &intVal
	}
	if internaltypes.IsDefined(plan.TrustedCertificateGroupId) {
		intVal := plan.TrustedCertificateGroupId.ValueInt64()
		addRequest.TrustedCertificateGroupId = &intVal
	}
	if internaltypes.IsDefined(plan.Secure) {
		boolVal := plan.Secure.ValueBool()
		addRequest.Secure = &boolVal
	}
	if internaltypes.IsDefined(plan.SendPaCookie) {
		boolVal := plan.SendPaCookie.ValueBool()
		addRequest.SendPaCookie = &boolVal
	}
	if internaltypes.IsDefined(plan.SkipHostnameVerification) {
		boolVal := plan.SkipHostnameVerification.ValueBool()
		addRequest.SkipHostnameVerification = &boolVal
	}
	if internaltypes.IsDefined(plan.UseProxy) {
		boolVal := plan.UseProxy.ValueBool()
		addRequest.UseProxy = &boolVal
	}
	if internaltypes.IsDefined(plan.UseTargetHostHeader) {
		boolVal := plan.UseTargetHostHeader.ValueBool()
		addRequest.UseTargetHostHeader = &boolVal
	}
	return nil
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSiteResponse(ctx context.Context, r *client.Site, state *siteResourceModel, expectedValues *siteResourceModel) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.AvailabilityProfileId = types.Int64Value(*r.AvailabilityProfileId)
	state.ExpectedHostname = internaltypes.StringTypeOrNil(r.ExpectedHostname, false)
	state.KeepAliveTimeout = types.Int64Value(*r.KeepAliveTimeout)
	state.LoadBalancingStrategyId = types.Int64Value(*r.LoadBalancingStrategyId)
	state.MaxConnections = types.Int64Value(*r.MaxConnections)
	state.MaxWebSocketConnections = types.Int64Value(*r.MaxWebSocketConnections)
	state.Name = types.StringValue(r.Name)
	state.Secure = internaltypes.BoolTypeOrNil(r.Secure)
	state.SendPaCookie = internaltypes.BoolTypeOrNil(r.SendPaCookie)
	state.SiteAuthenticatorIds = internaltypes.GetInt64Set(r.SiteAuthenticatorIds)
	state.SkipHostnameVerification = internaltypes.BoolTypeOrNil(r.SkipHostnameVerification)
	state.Targets = internaltypes.GetStringSet(r.Targets)
	state.TrustedCertificateGroupId = types.Int64Value(*r.TrustedCertificateGroupId)
	state.UseProxy = internaltypes.BoolTypeOrNil(r.UseProxy)
	state.UseTargetHostHeader = internaltypes.BoolTypeOrNil(r.UseTargetHostHeader)
}

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan siteResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var TargetsSlice []string
	plan.Targets.ElementsAs(ctx, &TargetsSlice, false)
	createSite := client.NewSite(plan.Name.ValueString(), TargetsSlice)
	err := addOptionalSiteFields(ctx, createSite, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Site", err.Error())
		return
	}
	requestJson, err := createSite.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateSite := r.apiClient.SitesApi.AddSite(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSite = apiCreateSite.Site(*createSite)
	siteResponse, httpResp, err := r.apiClient.SitesApi.AddSiteExecute(apiCreateSite)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Site", err, httpResp)
		return
	}
	responseJson, err := siteResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state siteResourceModel

	readSiteResponse(ctx, siteResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSite(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSite(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state siteResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSite, httpResp, err := apiClient.SitesApi.GetSite(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a Site", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadSite.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSiteResponse(ctx, apiReadSite, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSite(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSite(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state siteResourceModel
	req.State.Get(ctx, &state)
	var TargetsSlice []string
	plan.Targets.ElementsAs(ctx, &TargetsSlice, false)
	UpdateSite := apiClient.SitesApi.UpdateSite(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewSite(plan.Name.ValueString(), TargetsSlice)
	err := addOptionalSiteFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Site", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateSite = UpdateSite.Site(*CreateUpdateRequest)
	updateSiteResponse, httpResp, err := apiClient.SitesApi.UpdateSiteExecute(UpdateSite)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Site", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateSiteResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readSiteResponse(ctx, updateSiteResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteSite(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteSite(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.SitesApi.DeleteSite(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Site", err, httpResp)
		return
	}

}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
