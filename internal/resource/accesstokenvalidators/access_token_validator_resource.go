package accessTokenValidators

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
	_ resource.Resource                = &accessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &accessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &accessTokenValidatorResource{}
)

// EngineListenerResource is a helper function to simplify the provider implementation.
func AccessTokenValidatorResource() resource.Resource {
	return &accessTokenValidatorResource{}
}

// engineListenerResource is the resource implementation.
type accessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type accessTokenValidatorResourceModel struct {
	ClassName     types.String `tfsdk:"classname"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Configuration types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *accessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accessTokenValidatorResourceSchema(ctx, req, resp, false)
}

func accessTokenValidatorResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages Access Token Validator.",
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
			"classname": schema.StringAttribute{
				Required: true,
			},
			"configuration": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"description": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"path": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"subject_attribute_name": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"issuer": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"audience": schema.StringAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
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

func addOptionalAccessTokenValidatorFields(ctx context.Context, addRequest *client.AccessTokenValidator, plan accessTokenValidatorResourceModel) error {
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
func (r *accessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_token_validator"
}

func (r *accessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readAccessTokenValidatorResponse(ctx context.Context, r *client.AccessTokenValidator, state *accessTokenValidatorResourceModel, expectedValues *accessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.Name = types.StringValue(r.Name)
	state.ClassName = types.StringValue(r.ClassName)
	attrTypes := map[string]attr.Type{
		"audience":               basetypes.StringType{},
		"description":            basetypes.StringType{},
		"issuer":                 basetypes.StringType{},
		"path":                   basetypes.StringType{},
		"subject_attribute_name": basetypes.StringType{},
	}

	configValues := r.GetConfiguration()
	attrValues := map[string]attr.Value{
		"audience":               internaltypes.StringValueOrNull(configValues["audience"]),
		"description":            internaltypes.StringValueOrNull(configValues["description"]),
		"issuer":                 internaltypes.StringValueOrNull(configValues["issuer"]),
		"path":                   internaltypes.StringValueOrNull(configValues["path"]),
		"subject_attribute_name": internaltypes.StringValueOrNull(configValues["subjectAttributeName"]),
	}
	configuration := internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
	state.Configuration = configuration
}

func (r *accessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accessTokenValidatorResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createAccessTokenValidator := client.NewAccessTokenValidator(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalAccessTokenValidatorFields(ctx, createAccessTokenValidator, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Access Token Validator", err.Error())
		return
	}
	requestJson, err := createAccessTokenValidator.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiCreateAccessTokenValidator := r.apiClient.AccessTokenValidatorsApi.AddAccessTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAccessTokenValidator = apiCreateAccessTokenValidator.AccessTokenValidator(*createAccessTokenValidator)
	accessTokenValidatorResponse, httpResp, err := r.apiClient.AccessTokenValidatorsApi.AddAccessTokenValidatorExecute(apiCreateAccessTokenValidator)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating Access Token Validator", err, httpResp)
		return
	}
	responseJson, err := accessTokenValidatorResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accessTokenValidatorResourceModel

	readAccessTokenValidatorResponse(ctx, accessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAccessTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state accessTokenValidatorResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAccessTokenValidator, httpResp, err := apiClient.AccessTokenValidatorsApi.GetAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an Access Token Validator", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAccessTokenValidator.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAccessTokenValidatorResponse(ctx, apiReadAccessTokenValidator, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAccessTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan accessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state accessTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	UpdateAccessTokenValidator := apiClient.AccessTokenValidatorsApi.UpdateAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewAccessTokenValidator(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalAccessTokenValidatorFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for Access Token Validator", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateAccessTokenValidator = UpdateAccessTokenValidator.AccessTokenValidator(*CreateUpdateRequest)
	UpdateAccessTokenValidatorResponse, httpResp, err := apiClient.AccessTokenValidatorsApi.UpdateAccessTokenValidatorExecute(UpdateAccessTokenValidator)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Access Token Validator", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := UpdateAccessTokenValidatorResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readAccessTokenValidatorResponse(ctx, UpdateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *accessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAccessTokenValidator(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state accessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.AccessTokenValidatorsApi.DeleteAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting Access Token Validator", err, httpResp)
		return
	}
}

func (r *accessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
