package accessTokenValidators

import (
	"context"

	config "terraform-provider-pingaccess/internal/resource"
	internaltypes "terraform-provider-pingaccess/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
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

type accessTokenValidatorConfigModel struct {
	Description          types.String `tfsdk:"description"`
	Audience             types.String `tfsdk:"audience"`
	Issuer               types.String `tfsdk:"issuer"`
	Path                 types.String `tfsdk:"path"`
	SubjectAttributeName types.String `tfsdk:"subject_attribute_name"`
}

type accessTokenValidatorResourceModel struct {
	ClassName     types.String `tfsdk:"classname"`
	Id            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Configuration types.Object `tfsdk:"configuration"`
}

// GetSchema defines the schema for the resource.
func (r *accessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accessTokenValidatorResourceSchema(ctx, req, resp, false)
}

func accessTokenValidatorResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages Access token Validator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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

func addOptionalAccessTokenValidatorFields(ctx context.Context, addRequest *client.AccessTokenValidator, plan accessTokenValidatorResourceModel) diag.Diagnostics {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		intVal := plan.Id.ValueInt64()
		addRequest.Id = &intVal
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
	state.Id = types.Int64Value(*r.Id)
	state.Name = types.StringValue(r.Name)
	state.ClassName = types.StringValue(r.ClassName)
	attrTypes := map[string]attr.Type{
		"issuer":                 basetypes.StringType{},
		"path":                   basetypes.StringType{},
		"audience":               basetypes.StringType{},
		"subject_attribute_name": basetypes.StringType{},
		"description":            basetypes.StringType{},
	}

	configValues := r.GetConfiguration()
	attrValues := map[string]attr.Value{
		"description":            internaltypes.StringValueOrNull(configValues["description"]),
		"path":                   internaltypes.StringValueOrNull(configValues["path"]),
		"subject_attribute_name": internaltypes.StringValueOrNull(configValues["subjectAttributeName"]),
		"issuer":                 internaltypes.StringValueOrNull(configValues["issuer"]),
		"audience":               internaltypes.StringValueOrNull(configValues["audience"]),
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
	addOptionalAccessTokenValidatorFields(ctx, createAccessTokenValidator, plan)
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
	apiReadAccessTokenValidator, httpResp, err := apiClient.AccessTokenValidatorsApi.GetAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id)).Execute()

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
	UpdateAccessTokenValidator := apiClient.AccessTokenValidatorsApi.UpdateAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(plan.Id))
	CreateUpdateRequest := client.NewAccessTokenValidator(plan.ClassName.ValueString(), plan.Name.ValueString())
	addOptionalAccessTokenValidatorFields(ctx, CreateUpdateRequest, plan)
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateAccessTokenValidator = UpdateAccessTokenValidator.AccessTokenValidator(*CreateUpdateRequest)
	UpdateAccessTokenValidatorResponse, httpResp, err := apiClient.AccessTokenValidatorsApi.UpdateAccessTokenValidatorExecute(UpdateAccessTokenValidator)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating engine listener", err, httpResp)
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
	httpResp, err := apiClient.AccessTokenValidatorsApi.DeleteAccessTokenValidator(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting Access Token Validator", err, httpResp)
		return
	}
}

func (r *accessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
