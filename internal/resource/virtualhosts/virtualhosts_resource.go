package virtualhost

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &virtualhostResource{}
	_ resource.ResourceWithConfigure   = &virtualhostResource{}
	_ resource.ResourceWithImportState = &virtualhostResource{}
)

// VirtualHostResource is a helper function to simplify the provider implementation.
func VirtualHostResource() resource.Resource {
	return &virtualhostResource{}
}

// virtualhostResource is the resource implementation.
type virtualhostResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type virtualhostResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	AgentResourceCacheTTL     types.Int64  `tfsdk:"agent_resource_cache_ttl"`
	Host                      types.String `tfsdk:"host"`
	KeyPairId                 types.Int64  `tfsdk:"keypair_id"`
	Port                      types.Int64  `tfsdk:"port"`
	TrustedCertificateGroupId types.Int64  `tfsdk:"trusted_certificate_group_id"`
}

// GetSchema defines the schema for the resource.
func (r *virtualhostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	virtualhostResourceSchema(ctx, req, resp, false)
}

func virtualhostResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a VirtualHost.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_resource_cache_ttl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Required: true,
			},
			"keypair_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"trusted_certificate_group_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"port", "host"})
	}
	resp.Schema = schema
}
func addOptionalVirtualHostFields(ctx context.Context, addRequest *client.VirtualHost, plan virtualhostResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsDefined(plan.AgentResourceCacheTTL) {
		intVal := plan.AgentResourceCacheTTL.ValueInt64()
		addRequest.AgentResourceCacheTTL = &intVal
	}
	if internaltypes.IsDefined(plan.KeyPairId) {
		intVal := plan.KeyPairId.ValueInt64()
		addRequest.KeyPairId = &intVal
	}

	if internaltypes.IsDefined(plan.TrustedCertificateGroupId) {
		intVal := plan.TrustedCertificateGroupId.ValueInt64()
		addRequest.TrustedCertificateGroupId = &intVal
	}

	return nil
}

// Metadata returns the resource type name.
func (r *virtualhostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtualhosts"
}

func (r *virtualhostResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readVirtualHostResponse(ctx context.Context, r *client.VirtualHost, state *virtualhostResourceModel, expectedValues *virtualhostResourceModel) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.AgentResourceCacheTTL = types.Int64Value(int64(*r.AgentResourceCacheTTL))
	state.Host = types.StringValue(r.Host)
	state.KeyPairId = types.Int64Value(int64(*r.KeyPairId))
	state.Port = types.Int64Value(r.Port)
	state.TrustedCertificateGroupId = types.Int64Value(*r.TrustedCertificateGroupId)
}

func (r *virtualhostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualhostResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createVirtualHost := client.NewVirtualHost(plan.Host.ValueString(), plan.Port.ValueInt64())
	err := addOptionalVirtualHostFields(ctx, createVirtualHost, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for VirtualHost", err.Error())
		return
	}
	requestJson, err := createVirtualHost.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateVirtualHost := r.apiClient.VirtualhostsApi.AddVirtualHost(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateVirtualHost = apiCreateVirtualHost.VirtualHost(*createVirtualHost)
	listenerResponse, httpResp, err := r.apiClient.VirtualhostsApi.AddVirtualHostExecute(apiCreateVirtualHost)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the VirtualHost", err, httpResp)
		return
	}
	responseJson, err := listenerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualhostResourceModel

	readVirtualHostResponse(ctx, listenerResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualhostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVirtualHost(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readVirtualHost(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state virtualhostResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadVirtualHost, httpResp, err := apiClient.VirtualhostsApi.GetVirtualHost(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a VirtualHost", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadVirtualHost.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVirtualHostResponse(ctx, apiReadVirtualHost, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *virtualhostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVirtualHost(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVirtualHost(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan virtualhostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state virtualhostResourceModel
	req.State.Get(ctx, &state)
	UpdateVirtualHost := apiClient.VirtualhostsApi.UpdateVirtualHost(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewVirtualHost(plan.Host.ValueString(), plan.Port.ValueInt64())
	err := addOptionalVirtualHostFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for VirtualHost", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateVirtualHost = UpdateVirtualHost.VirtualHost(*CreateUpdateRequest)
	updateVirtualHostResponse, httpResp, err := apiClient.VirtualhostsApi.UpdateVirtualHostExecute(UpdateVirtualHost)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating VirtualHost", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateVirtualHostResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readVirtualHostResponse(ctx, updateVirtualHostResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *virtualhostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteVirtualHost(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteVirtualHost(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state virtualhostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.VirtualhostsApi.DeleteVirtualHost(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a VirtualHost", err, httpResp)
		return
	}

}

func (r *virtualhostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
