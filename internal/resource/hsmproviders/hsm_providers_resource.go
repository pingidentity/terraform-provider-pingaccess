package hsmprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &hsmProviderResource{}
	_ resource.ResourceWithConfigure   = &hsmProviderResource{}
	_ resource.ResourceWithImportState = &hsmProviderResource{}
)

// HsmProviderResource is a helper function to simplify the provider implementation.
func HsmProviderResource() resource.Resource {
	return &hsmProviderResource{}
}

// hsmProviderResource is the resource implementation.
type hsmProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type hsmProviderResourceModel struct {
	Id            types.String `tfsdk:"id"`
	ClassName     types.String `tfsdk:"classname"`
	Configuration types.Object `tfsdk:"configuration"`
	Name          types.String `tfsdk:"name"`
}

// GetSchema defines the schema for the resource.
func (r *hsmProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	hsmProviderResourceSchema(ctx, req, resp, false)
}

func hsmProviderResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a HsmProvider.",
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
					"user": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"partition": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"slot_id": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"library": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"classname", "name", "configuration"})
	}
	resp.Schema = schema
}

func addOptionalHsmProviderFields(ctx context.Context, addRequest *client.HsmProvider, plan hsmProviderResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsNonEmptyObj(plan.Configuration) {

		addRequest.Configuration = internaltypes.ObjValuesToClientMap(plan.Configuration)
	}
	return nil
}

// Metadata returns the resource type name.
func (r *hsmProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hsm_providers"
}

func (r *hsmProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *hsmProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model hsmProviderResourceModel
	req.Plan.Get(ctx, &model)
	switch model.ClassName.ValueString() {
	case "com.pingidentity.pa.hsm.pkcs11.plugin.PKCS11HsmProvider":
		if internaltypes.IsNonEmptyString(model.Configuration.Attributes()["user"].(basetypes.StringValue)) {
			resp.Diagnostics.AddError("Attribute 'user' not supported by SafenetHSM Provider", "Required attributes are slot_id,password and library")
		}
		if internaltypes.IsNonEmptyString(model.Configuration.Attributes()["partition"].(basetypes.StringValue)) {
			resp.Diagnostics.AddError("Attribute 'partition' not supported by SafenetHSM Provider", "Required attributes are slot_id,password and library")
		}
	case "com.pingidentity.pa.hsm.cloudhsm.plugin.AwsCloudHsmProvider":
		if internaltypes.IsNonEmptyString(model.Configuration.Attributes()["slot_id"].(basetypes.StringValue)) {
			resp.Diagnostics.AddError("Attribute 'slot_id' not supported by AwsHSM Provider", "Required attributes are user,password and partition")
		}
		if internaltypes.IsNonEmptyString(model.Configuration.Attributes()["library"].(basetypes.StringValue)) {
			resp.Diagnostics.AddError("Attribute 'library' not supported by AwsHSM Provider", "Required attributes are user,password and partition")
		}
	}
}

func readHsmProviderResponse(ctx context.Context, r *client.HsmProvider, state *hsmProviderResourceModel, expectedValues *hsmProviderResourceModel, diagnostics *diag.Diagnostics, createPlan basetypes.ObjectValue) {
	passwordConfiguration := *internaltypes.ObjValuesToClientMap(createPlan)
	state.Id = internaltypes.StringValueOrNull(internaltypes.Int64PointerToString(*r.Id))
	state.ClassName = types.StringValue(r.ClassName)
	state.Name = types.StringValue(r.Name)
	attrTypes := map[string]attr.Type{
		"user":      basetypes.StringType{},
		"password":  basetypes.StringType{},
		"partition": basetypes.StringType{},
		"slot_id":   basetypes.StringType{},
		"library":   basetypes.StringType{},
	}

	configValues := r.GetConfiguration()
	attrValues := map[string]attr.Value{
		"user":      internaltypes.StringValueOrNull(configValues["user"]),
		"password":  internaltypes.StringValueOrNull(passwordConfiguration["password"]),
		"partition": internaltypes.StringValueOrNull(configValues["partition"]),
		"slot_id":   internaltypes.StringValueOrNull(configValues["slotId"]),
		"library":   internaltypes.StringValueOrNull(configValues["library"]),
	}
	configuration := internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
	state.Configuration = configuration

}

func (r *hsmProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hsmProviderResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// passplain := plan.Configuration.Attributes()
	createHsmProvider := client.NewHsmProvider(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalHsmProviderFields(ctx, createHsmProvider, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Hsm Provider", err.Error())
		return
	}
	requestJson, err := createHsmProvider.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateHsmProvider := r.apiClient.HsmProvidersApi.AddHsmProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateHsmProvider = apiCreateHsmProvider.HsmProvider(*createHsmProvider)
	hsmResponse, httpResp, err := r.apiClient.HsmProvidersApi.AddHsmProviderExecute(apiCreateHsmProvider)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the HsmProvider", err, httpResp)
		return
	}
	responseJson, err := hsmResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state hsmProviderResourceModel

	readHsmProviderResponse(ctx, hsmResponse, &state, &plan, &resp.Diagnostics, plan.Configuration)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hsmProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHsmProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readHsmProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state hsmProviderResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadHsmProvider, httpResp, err := apiClient.HsmProvidersApi.GetHsmProvider(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an HsmProvider", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadHsmProvider.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHsmProviderResponse(ctx, apiReadHsmProvider, &state, &state, &resp.Diagnostics, state.Configuration)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *hsmProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHsmProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateHsmProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan hsmProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state hsmProviderResourceModel
	req.State.Get(ctx, &state)
	UpdateHsmProvider := apiClient.HsmProvidersApi.UpdateHsmProvider(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewHsmProvider(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalHsmProviderFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for Hsm Provider", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateHsmProvider = UpdateHsmProvider.HsmProvider(*CreateUpdateRequest)
	updateHsmProviderResponse, httpResp, err := apiClient.HsmProvidersApi.UpdateHsmProviderExecute(UpdateHsmProvider)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating HsmProvider", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateHsmProviderResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readHsmProviderResponse(ctx, updateHsmProviderResponse, &state, &plan, &resp.Diagnostics, plan.Configuration)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *hsmProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteHsmProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteHsmProvider(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state hsmProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	switch state.ClassName.ValueString() {
	case "com.pingidentity.pa.hsm.pkcs11.plugin.PKCS11HsmProvider":
		httpResp, err := apiClient.HsmProvidersApi.DeleteHsmProvider(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting Safenet HsmProvider", err, httpResp)
			return
		}

	case "com.pingidentity.pa.hsm.cloudhsm.plugin.AwsCloudHsmProvider":
		httpResp, err := apiClient.HsmProvidersApi.DeleteHsmProvider(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting Aws HsmProvider.\nPlease remove the custom aws jar file from the deploy folder before removing the AwsProvider.", err, httpResp)
			return
		}
	}
}

func (r *hsmProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
