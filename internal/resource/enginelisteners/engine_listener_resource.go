package engineListener

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
	_ resource.Resource                = &engineListenerResource{}
	_ resource.ResourceWithConfigure   = &engineListenerResource{}
	_ resource.ResourceWithImportState = &engineListenerResource{}
)

// EngineListenerResource is a helper function to simplify the provider implementation.
func EngineListenerResource() resource.Resource {
	return &engineListenerResource{}
}

// engineListenerResource is the resource implementation.
type engineListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type engineListenerResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Port                      types.Int64  `tfsdk:"port"`
	Secure                    types.Bool   `tfsdk:"secure"`
	TrustedCertificateGroupId types.Int64  `tfsdk:"trusted_certificate_group_id"`
}

// GetSchema defines the schema for the resource.
func (r *engineListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	enginelistenerResourceSchema(ctx, req, resp, false)
}

func enginelistenerResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages an Engine Listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"port": schema.Int64Attribute{
				Required: true,
			},
			"secure": schema.BoolAttribute{
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
		},
	}

	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "port"})
	}
	resp.Schema = schema
}
func addOptionalEngineListenerFields(ctx context.Context, addRequest *client.EngineListener, plan engineListenerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsDefined(plan.Secure) {
		boolVal := plan.Secure.ValueBool()
		addRequest.Secure = &boolVal
	}
	if internaltypes.IsDefined(plan.TrustedCertificateGroupId) {
		intVal := plan.TrustedCertificateGroupId.ValueInt64()
		addRequest.TrustedCertificateGroupId = &intVal
	}
	return nil
}

// Metadata returns the resource type name.
func (r *engineListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine_listener"
}

func (r *engineListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readEngineListenerResponse(ctx context.Context, r *client.EngineListener, state *engineListenerResourceModel, expectedValues *engineListenerResourceModel) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.Name = types.StringValue(r.Name)
	state.Port = types.Int64Value(int64(r.Port))
	state.Secure = internaltypes.BoolTypeOrNil(r.Secure)
	state.TrustedCertificateGroupId = types.Int64Value(*r.TrustedCertificateGroupId)
}

func (r *engineListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan engineListenerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createListener := client.NewEngineListener(plan.Name.ValueString(), plan.Port.ValueInt64())
	err := addOptionalEngineListenerFields(ctx, createListener, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Engine Listener", err.Error())
		return
	}
	requestJson, err := createListener.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateListener := r.apiClient.EngineListenersApi.AddEngineListener(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateListener = apiCreateListener.EngineListener(*createListener)
	listenerResponse, httpResp, err := r.apiClient.EngineListenersApi.AddEngineListenerExecute(apiCreateListener)

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the engine listener", err, httpResp)
		return
	}
	responseJson, err := listenerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state engineListenerResourceModel

	readEngineListenerResponse(ctx, listenerResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *engineListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readEngineListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readEngineListener(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state engineListenerResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadListener, httpResp, err := apiClient.EngineListenersApi.GetEngineListener(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an engine listener", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadListener.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readEngineListenerResponse(ctx, apiReadListener, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engineListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateEngineListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateEngineListener(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan engineListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state engineListenerResourceModel
	req.State.Get(ctx, &state)
	UpdateListener := apiClient.EngineListenersApi.UpdateEngineListener(config.ProviderBasicAuthContext(ctx, providerConfig), (plan.Id.ValueString()))
	CreateUpdateRequest := client.NewEngineListener(plan.Name.ValueString(), plan.Port.ValueInt64())
	err := addOptionalEngineListenerFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for AcmeServer", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateListener = UpdateListener.EngineListener(*CreateUpdateRequest)
	updateListenerResponse, httpResp, err := apiClient.EngineListenersApi.UpdateEngineListenerExecute(UpdateListener)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating engine listener", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateListenerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readEngineListenerResponse(ctx, updateListenerResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *engineListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteEngineListener(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteEngineListener(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state engineListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.EngineListenersApi.DeleteEngineListener(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an engine listener", err, httpResp)
		return
	}

}

func (r *engineListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
