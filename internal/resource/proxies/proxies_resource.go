package proxie

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &proxieResource{}
	_ resource.ResourceWithConfigure   = &proxieResource{}
	_ resource.ResourceWithImportState = &proxieResource{}
)

// HttpClientProxyResource is a helper function to simplify the provider implementation.
func HttpClientProxyResource() resource.Resource {
	return &proxieResource{}
}

// proxieResource is the resource implementation.
type proxieResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type proxieResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Description            types.String `tfsdk:"description"`
	Host                   types.String `tfsdk:"host"`
	Name                   types.String `tfsdk:"name"`
	Password               types.Object `tfsdk:"password"`
	Port                   types.Int64  `tfsdk:"port"`
	RequiresAuthentication types.Bool   `tfsdk:"requires_authentication"`
	Username               types.String `tfsdk:"username"`
}

// GetSchema defines the schema for the resource.
func (r *proxieResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	proxieResourceSchema(ctx, req, resp, false)
}

func proxieResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a HttpClientProxy.",
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
			"host": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.Int64Attribute{
				Required: true,
			},
			"password": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Sensitive: true,
						Required:  true,
					},
				},
			},
			"requires_authentication": schema.BoolAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "host", "port"})
	}
	resp.Schema = schema
}
func addOptionalHttpClientProxyFields(ctx context.Context, addRequest *client.HttpClientProxy, plan proxieResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}
	if internaltypes.IsDefined(plan.RequiresAuthentication) {
		boolVal := plan.RequiresAuthentication.ValueBool()
		addRequest.RequiresAuthentication = &boolVal
	}
	if internaltypes.IsDefined(plan.Username) {
		stringVal := plan.Username.ValueString()
		addRequest.Username = &stringVal
	}
	if internaltypes.IsDefined(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	addRequest.Password = client.NewHiddenField()
	if internaltypes.IsDefined(plan.Password) {
		pass := plan.Password.Attributes()
		passValue := pass["value"]
		if !passValue.IsNull() || !passValue.IsUnknown() {
			addRequest.Password.Value = internaltypes.InterfaceStringPointerValue(internaltypes.ConvertToPrimitive(pass["value"]))
		}
	}
	return nil
}

// Metadata returns the resource type name.
func (r *proxieResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy"
}

func (r *proxieResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readHttpClientProxyResponse(ctx context.Context, r *client.HttpClientProxy, state *proxieResourceModel, expectedValues *proxieResourceModel, createPlan basetypes.ObjectValue, addLastEncrypted bool) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Host = types.StringValue(r.Host)
	state.Name = types.StringValue(r.Name)

	passwordConfiguration := *internaltypes.ObjValuesToClientMap(createPlan)
	attrTypes := map[string]attr.Type{
		"value": basetypes.StringType{},
	}
	attrValues := map[string]attr.Value{
		"value": internaltypes.StringValueOrNull(passwordConfiguration["value"]),
	}

	password := internaltypes.MaptoObjValue(attrTypes, attrValues, diag.Diagnostics{})
	state.Password = password
	state.Port = types.Int64Value(r.Port)
	state.RequiresAuthentication = internaltypes.BoolTypeOrNil(r.RequiresAuthentication)
	state.Username = internaltypes.StringTypeOrNil(r.Username, false)

}

func (r *proxieResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan proxieResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createHttpClientProxy := client.NewHttpClientProxy(plan.Name.ValueString(), plan.Host.ValueString(), plan.Port.ValueInt64())
	err := addOptionalHttpClientProxyFields(ctx, createHttpClientProxy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for HttpClientProxy", err.Error())
		return
	}
	requestJson, err := createHttpClientProxy.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateHttpClientProxy := r.apiClient.ProxiesApi.AddProxy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateHttpClientProxy = apiCreateHttpClientProxy.Proxy(*createHttpClientProxy)
	proxieResponse, httpResp, err := r.apiClient.ProxiesApi.AddProxyExecute(apiCreateHttpClientProxy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the HttpClientProxy", err, httpResp)
		return
	}
	responseJson, err := proxieResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state proxieResourceModel

	readHttpClientProxyResponse(ctx, proxieResponse, &state, &plan, plan.Password, true)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *proxieResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpClientProxy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readHttpClientProxy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state proxieResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadHttpClientProxy, httpResp, err := apiClient.ProxiesApi.GetProxy(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a HttpClientProxy", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadHttpClientProxy.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpClientProxyResponse(ctx, apiReadHttpClientProxy, &state, &state, state.Password, false)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *proxieResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpClientProxy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateHttpClientProxy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan proxieResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state proxieResourceModel
	req.State.Get(ctx, &state)

	UpdateHttpClientProxy := apiClient.ProxiesApi.UpdateProxy(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewHttpClientProxy(plan.Name.ValueString(), plan.Host.ValueString(), plan.Port.ValueInt64())
	err := addOptionalHttpClientProxyFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for HttpClientProxy", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateHttpClientProxy = UpdateHttpClientProxy.Proxy(*CreateUpdateRequest)
	updateHttpClientProxyResponse, httpResp, err := apiClient.ProxiesApi.UpdateProxyExecute(UpdateHttpClientProxy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating HttpClientProxy", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateHttpClientProxyResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readHttpClientProxyResponse(ctx, updateHttpClientProxyResponse, &state, &plan, plan.Password, true)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *proxieResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteHttpClientProxy(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteHttpClientProxy(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state proxieResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.ProxiesApi.DeleteProxy(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a HttpClientProxy", err, httpResp)
		return
	}

}

func (r *proxieResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
