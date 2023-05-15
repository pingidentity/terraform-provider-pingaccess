package identitymappings

import (
	"context"

	config "github.com/pingidentity/terraform-provider-pingaccess/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"

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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &identityMappingResource{}
	_ resource.ResourceWithConfigure   = &identityMappingResource{}
	_ resource.ResourceWithImportState = &identityMappingResource{}
)

// IdentityMappingResource is a helper function to simplify the provider implementation.
func IdentityMappingResource() resource.Resource {
	return &identityMappingResource{}
}

// identityMappingResource is the resource implementation.
type identityMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type identityMappingResourceModel struct {
	Id            types.String `tfsdk:"id"`
	ClassName     types.String `tfsdk:"classname"`
	Configuration types.Object `tfsdk:"configuration"`
	Name          types.String `tfsdk:"name"`
}

// GetSchema defines the schema for the resource.
func (r *identityMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	identityMappingResourceSchema(ctx, req, resp, false)
}

func identityMappingResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a IdentityMapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"classname": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configuration": schema.SingleNestedAttribute{
				Optional: true,
				// Default: ,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"attribute_header_mappings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_name": schema.StringAttribute{
									Optional: true,
								},
								"header_name": schema.StringAttribute{
									Optional: true,
								},
								"subject": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"attribute_mappings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"jwt_claim_name": schema.StringAttribute{
									Optional: true,
								},
								"user_attribute_name": schema.StringAttribute{
									Optional: true,
								},
								"subject": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"audience": schema.StringAttribute{
						Optional: true,
					},
					"cache_jwt": schema.BoolAttribute{
						Optional: true,
					},
					"client_certificate_jwt_claim_name": schema.StringAttribute{
						Optional: true,
					},
					"exclusion_list": schema.BoolAttribute{
						Optional: true,
					},
					"exclusion_list_attributes": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"exclusion_list_subject": schema.StringAttribute{
						Optional: true,
					},
					"header_client_certificate_mappings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"header_name": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"header_name": schema.StringAttribute{
						Optional: true,
					},
					"header_name_prefix": schema.StringAttribute{
						Optional: true,
					},
					"map_as_bearer_token": schema.BoolAttribute{
						Optional: true,
					},
					"max_depth": schema.Float64Attribute{
						Optional: true,
					},
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "classname", "configuration"})
	}
	resp.Schema = schema
}
func addOptionalIdentityMappingFields(ctx context.Context, addRequest *client.IdentityMapping, plan identityMappingResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsNonEmptyObj(plan.Configuration) {
		converted := internaltypes.ConvertToPrimitive(plan.Configuration)
		mapConverted := converted.(map[string]interface{})
		addRequest.Configuration = mapConverted
	}
	return nil
}

// Metadata returns the resource type name.
func (r *identityMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_mappings"
}

func (r *identityMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *identityMappingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model identityMappingResourceModel
	req.Plan.Get(ctx, &model)
	configAttrKeys := model.Configuration.Attributes()
	switch model.ClassName.ValueString() {
	case "com.pingidentity.pa.identitymappings.HeaderIdentityMapping":
		himReqConfigProps := []string{"header_name_prefix", "exclusion_list_attributes", "attribute_header_mappings", "header_client_certificate_mappings"}
		for configAttrKey := range configAttrKeys {
			checkKeyMatch := internaltypes.StringSliceContains(himReqConfigProps, configAttrKey)
			hasConfigValue := internaltypes.IsDefined(configAttrKeys[configAttrKey])
			if !checkKeyMatch && hasConfigValue {
				resp.Diagnostics.AddError("Attribute "+configAttrKey+" not supported for Header Identity Mapping!", "")
			}
		}
	case "com.pingidentity.pa.identitymappings.JwtIdentityMapping":
		jwtReqConfigProps := []string{"map_as_bearer_token", "header_name", "audience", "exclusion_list_attributes", "exclusion_list_subject", "exclusion_list", "attribute_mappings", "cache_jwt", "client_certificate_jwt_claim_name", "max_depth"}
		for configAttrKey := range configAttrKeys {
			checkKeyMatch := internaltypes.StringSliceContains(jwtReqConfigProps, configAttrKey)
			hasConfigValue := internaltypes.IsDefined(configAttrKeys[configAttrKey])
			if !checkKeyMatch && hasConfigValue {
				resp.Diagnostics.AddError("Attribute "+configAttrKey+" not supported for Jwt Identity Mapping!", "")
			}
		}
	}
}

func readIdentityMappingResponse(ctx context.Context, r *client.IdentityMapping, state *identityMappingResourceModel, expectedValues *identityMappingResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.ClassName = types.StringValue(r.ClassName)
	state.Name = types.StringValue(r.Name)
	configValues := r.GetConfiguration()
	className = types.StringValue(r.ClassName)
	attrTypes := map[string]attr.Type{
		"attribute_header_mappings": basetypes.SetType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"jwt_claim_name":      basetypes.StringType{},
				"subject":             basetypes.BoolType{},
				"user_attribute_name": basetypes.StringType{},
			},
		}},
		"attribute_mappings": basetypes.SetType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"jwt_claim_name":      basetypes.StringType{},
				"subject":             basetypes.BoolType{},
				"user_attribute_name": basetypes.StringType{},
			},
		}},
		"audience":                          basetypes.StringType{},
		"cache_jwt":                         basetypes.BoolType{},
		"client_certificate_jwt_claim_name": basetypes.StringType{},
		"exclusion_list":                    basetypes.BoolType{},
		"exclusion_list_attributes":         basetypes.SetType{ElemType: types.StringType},
		"exclusion_list_subject":            basetypes.StringType{},
		"header_client_certificate_mappings": basetypes.SetType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"jwt_claim_name":      basetypes.StringType{},
				"subject":             basetypes.BoolType{},
				"user_attribute_name": basetypes.StringType{},
			},
		}},
		"header_name":         basetypes.StringType{},
		"header_name_prefix":  basetypes.StringType{},
		"map_as_bearer_token": basetypes.BoolType{},
		"max_depth":           basetypes.Float64Type{},
	}
	attrMappingsTypes := map[string]attr.Type{
		"jwt_claim_name":      basetypes.StringType{},
		"subject":             basetypes.BoolType{},
		"user_attribute_name": basetypes.StringType{},
	}
	attrHeaderMappingsTypes := map[string]attr.Type{
		"attribute_name": basetypes.StringType{},
		"header_name":    basetypes.StringType{},
		"subject":        basetypes.BoolType{},
	}
	headerClientCertificateMappings := map[string]attr.Type{
		"header_name": basetypes.StringType{},
	}

	attributeMappingsList, ok := configValues["attributeMappings"].([]interface{})
	if ok {
		attrValueSlice := []attr.Value{}
		for _, attrValuesMap := range attributeMappingsList {
			attrValuesMaptoInterface := attrValuesMap.(map[string]interface{})
			finalValues := map[string]attr.Value{}
			finalValues["subject"] = internaltypes.InterfaceBoolTypeOrNull(attrValuesMaptoInterface["subject"])
			finalValues["user_attribute_name"] = internaltypes.StringValueOrNull(attrValuesMaptoInterface["userAttributeName"])
			finalValues["jwt_claim_name"] = internaltypes.StringValueOrNull(attrValuesMaptoInterface["jwtClaimName"])
			attributeMappingObjects := internaltypes.MaptoObjValue(attrMappingsTypes, finalValues, *diagnostics)
			attrValueSlice = append(attrValueSlice, attributeMappingObjects)
		}
		attributeMappings, _ := types.SetValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"subject":             basetypes.BoolType{},
				"user_attribute_name": basetypes.StringType{},
				"jwt_claim_name":      basetypes.StringType{},
			},
		}, attrValueSlice)
		attrValues := map[string]attr.Value{
			"attribute_header_mappings":          internaltypes.InterfaceStringSetOrNil(configValues["attributeHeaderMappings"]),
			"attribute_mappings":                 attributeMappings,
			"audience":                           internaltypes.StringValueOrNull(configValues["audience"]),
			"cache_jwt":                          internaltypes.InterfaceBoolTypeOrNull(configValues["cacheJwt"]),
			"client_certificate_jwt_claim_name":  internaltypes.StringValueOrNull(configValues["clientCertificateJwtClaimName"]),
			"exclusion_list":                     internaltypes.InterfaceBoolTypeOrNull(configValues["exclusionList"]),
			"exclusion_list_attributes":          internaltypes.InterfaceStringSetOrNil(configValues["exclusionListAttributes"]),
			"exclusion_list_subject":             internaltypes.StringValueOrNull(configValues["exclusionListSubject"]),
			"header_client_certificate_mappings": internaltypes.GetNestedInterfaceKey(configValues["headerClientCertificateMappings"], "headerName"),
			"header_name":                        internaltypes.StringValueOrNull(configValues["headerName"]),
			"header_name_prefix":                 internaltypes.StringValueOrNull(configValues["headerNamePrefix"]),
			"map_as_bearer_token":                internaltypes.InterfaceBoolTypeOrNull(configValues["mapAsBearerToken"]),
			"max_depth":                          internaltypes.InterfaceFloat64TypeOrNull(configValues["maxDepth"]),
		}
		configuration := internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
		state.Configuration = configuration
	} else {
		attrValueSlice := []attr.Value{}
		for i := 0; i <= 1; i++ {
			finalValues := map[string]attr.Value{}
			finalValues["subject"] = types.BoolNull()
			finalValues["user_attribute_name"] = types.BoolNull()
			finalValues["jwt_claim_name"] = types.BoolNull()
			attributeMappingObjects := internaltypes.MaptoObjValue(attrMappingsTypes, finalValues, *diagnostics)
			attrValueSlice = append(attrValueSlice, attributeMappingObjects)
		}
		attributeMappings, _ := types.SetValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"subject":             basetypes.BoolType{},
				"user_attribute_name": basetypes.StringType{},
				"jwt_claim_name":      basetypes.StringType{},
			},
		}, attrValueSlice)

		attrValues := map[string]attr.Value{
			"attribute_header_mappings":          internaltypes.InterfaceStringSetOrNil(configValues["attributeHeaderMappings"]),
			"attribute_mappings":                 attributeMappings,
			"audience":                           internaltypes.StringValueOrNull(configValues["audience"]),
			"cache_jwt":                          internaltypes.InterfaceBoolTypeOrNull(configValues["cacheJwt"]),
			"client_certificate_jwt_claim_name":  internaltypes.StringValueOrNull(configValues["clientCertificateJwtClaimName"]),
			"exclusion_list":                     internaltypes.InterfaceBoolTypeOrNull(configValues["exclusionList"]),
			"exclusion_list_attributes":          internaltypes.InterfaceStringSetOrNil(configValues["exclusionListAttributes"]),
			"exclusion_list_subject":             internaltypes.StringValueOrNull(configValues["exclusionListSubject"]),
			"header_client_certificate_mappings": internaltypes.GetNestedInterfaceKey(configValues["headerClientCertificateMappings"], "headerName"),
			"header_name":                        internaltypes.StringValueOrNull(configValues["headerName"]),
			"header_name_prefix":                 internaltypes.StringValueOrNull(configValues["headerNamePrefix"]),
			"map_as_bearer_token":                internaltypes.InterfaceBoolTypeOrNull(configValues["mapAsBearerToken"]),
			"max_depth":                          internaltypes.InterfaceFloat64TypeOrNull(configValues["maxDepth"]),
		}
		state.Configuration = internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
	}
}

func (r *identityMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identityMappingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdentityMapping := client.NewIdentityMapping(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalIdentityMappingFields(ctx, createIdentityMapping, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdentityMapping", err.Error())
		return
	}
	requestJson, err := createIdentityMapping.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateIdentityMapping := r.apiClient.IdentityMappingsApi.AddIdentityMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdentityMapping = apiCreateIdentityMapping.IdentityMappings(*createIdentityMapping)
	identityMappingResponse, httpResp, err := r.apiClient.IdentityMappingsApi.AddIdentityMappingExecute(apiCreateIdentityMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdentityMapping", err, httpResp)
		return
	}
	responseJson, err := identityMappingResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMappingResourceModel

	readIdentityMappingResponse(ctx, identityMappingResponse, &state, &plan, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdentityMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdentityMapping(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state identityMappingResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdentityMapping, httpResp, err := apiClient.IdentityMappingsApi.GetIdentityMapping(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a IdentityMapping", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadIdentityMapping.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdentityMappingResponse(ctx, apiReadIdentityMapping, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *identityMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdentityMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdentityMapping(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan identityMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state identityMappingResourceModel
	req.State.Get(ctx, &state)
	UpdateIdentityMapping := apiClient.IdentityMappingsApi.UpdateIdentityMapping(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString())
	CreateUpdateRequest := client.NewIdentityMapping(plan.ClassName.ValueString(), plan.Name.ValueString())
	err := addOptionalIdentityMappingFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdentityMapping", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateIdentityMapping = UpdateIdentityMapping.IdentityMappings(*CreateUpdateRequest)
	updateIdentityMappingResponse, httpResp, err := apiClient.IdentityMappingsApi.UpdateIdentityMappingExecute(UpdateIdentityMapping)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdentityMapping", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateIdentityMappingResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readIdentityMappingResponse(ctx, updateIdentityMappingResponse, &state, &plan, &resp.Diagnostics)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *identityMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteIdentityMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteIdentityMapping(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state identityMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.IdentityMappingsApi.DeleteIdentityMapping(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a IdentityMapping", err, httpResp)
		return
	}

}

func (r *identityMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
