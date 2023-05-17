package authnReqList

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &authnReqListResource{}
	_ resource.ResourceWithConfigure   = &authnReqListResource{}
	_ resource.ResourceWithImportState = &authnReqListResource{}
)

// AuthnReqListResource is a helper function to simplify the provider implementation.
func AuthnReqListResource() resource.Resource {
	return &authnReqListResource{}
}

// authnReqListResource is the resource implementation.
type authnReqListResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authnReqListResourceModel struct {
	Id        types.String `tfsdk:"id"`
	AuthnReqs types.Set    `tfsdk:"authn_reqs"`
	Name      types.String `tfsdk:"name"`
}

// GetSchema defines the schema for the resource.
func (r *authnReqListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	authnReqListResourceSchema(ctx, req, resp, false)
}

func authnReqListResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a AuthnReqList.",
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
			"authn_reqs": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "authn_reqs"})
	}
	resp.Schema = schema
}
func addOptionalAuthnReqListFields(ctx context.Context, addRequest *client.AuthnReqList, plan authnReqListResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}

	return nil
}

// Metadata returns the resource type name.
func (r *authnReqListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authn_req_lists"
}

func (r *authnReqListResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthnReqListResponse(ctx context.Context, r *client.AuthnReqList, state *authnReqListResourceModel, expectedValues *authnReqListResourceModel) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.AuthnReqs = internaltypes.GetStringSet(r.AuthnReqs)
	state.Name = types.StringValue(r.Name)
}

func (r *authnReqListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authnReqListResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var authnReqs []string
	plan.AuthnReqs.ElementsAs(ctx, &authnReqs, false)
	createAuthnReqList := client.NewAuthnReqList(plan.Name.ValueString(), authnReqs)
	err := addOptionalAuthnReqListFields(ctx, createAuthnReqList, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthnReqList", err.Error())
		return
	}
	requestJson, err := createAuthnReqList.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateAuthnReqList := r.apiClient.AuthnReqListsApi.AddAuthnReqList(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthnReqList = apiCreateAuthnReqList.AuthnReqList(*createAuthnReqList)
	authnReqListResponse, httpResp, err := r.apiClient.AuthnReqListsApi.AddAuthnReqListExecute(apiCreateAuthnReqList)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the AuthnReqList", err, httpResp)
		return
	}
	responseJson, err := authnReqListResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state authnReqListResourceModel

	readAuthnReqListResponse(ctx, authnReqListResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *authnReqListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAuthnReqList(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAuthnReqList(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state authnReqListResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthnReqList, httpResp, err := apiClient.AuthnReqListsApi.GetAuthnReqList(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a AuthnReqList", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAuthnReqList.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthnReqListResponse(ctx, apiReadAuthnReqList, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authnReqListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAuthnReqList(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAuthnReqList(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan authnReqListResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authnReqListResourceModel
	req.State.Get(ctx, &state)
	var authnReqs []string
	plan.AuthnReqs.ElementsAs(ctx, &authnReqs, false)
	UpdateAuthnReqList := apiClient.AuthnReqListsApi.UpdateAuthnReqList(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewAuthnReqList(plan.Name.ValueString(), authnReqs)
	err := addOptionalAuthnReqListFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthnReqList", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateAuthnReqList = UpdateAuthnReqList.AuthnReqList(*CreateUpdateRequest)
	updateAuthnReqListResponse, httpResp, err := apiClient.AuthnReqListsApi.UpdateAuthnReqListExecute(UpdateAuthnReqList)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AuthnReqList", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAuthnReqListResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readAuthnReqListResponse(ctx, updateAuthnReqListResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *authnReqListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAuthnReqList(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAuthnReqList(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state authnReqListResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.AuthnReqListsApi.DeleteAuthnReqList(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a AuthnReqList", err, httpResp)
		return
	}

}

func (r *authnReqListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
