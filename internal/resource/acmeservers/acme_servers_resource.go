package acmeservers

import (
	"context"

	config "github.com/pingidentity/terraform-provider-pingaccess/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingaccess-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &acmeserversResource{}
	_ resource.ResourceWithConfigure   = &acmeserversResource{}
	_ resource.ResourceWithImportState = &acmeserversResource{}
)

// AcmeServerResource is a helper function to simplify the provider implementation.
func AcmeServerResource() resource.Resource {
	return &acmeserversResource{}
}

// acmeserversResource is the resource implementation.
type acmeserversResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type acmeserversResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Url  types.String `tfsdk:"url"`
}

// GetSchema defines the schema for the resource.
func (r *acmeserversResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	acmeserversResourceSchema(ctx, req, resp, false)
}

func acmeserversResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a AcmeServer.",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "url"})
	}
	resp.Schema = schema
}
func addOptionalAcmeServerFields(ctx context.Context, addRequest *client.AcmeServer, plan acmeserversResourceModel) error {
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}

	return nil
}

// Metadata returns the resource type name.
func (r *acmeserversResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acme_servers"
}

func (r *acmeserversResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAcmeServerResponse(ctx context.Context, r *client.AcmeServer, state *acmeserversResourceModel, expectedValues *acmeserversResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(*r.Id)
	state.Name = types.StringValue(r.Name)
	state.Url = types.StringValue(r.Url)
}

func (r *acmeserversResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan acmeserversResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAcmeServer := client.NewAcmeServer(plan.Name.ValueString(), plan.Url.ValueString())
	addOptionalAcmeServerFields(ctx, createAcmeServer, plan)
	requestJson, err := createAcmeServer.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateAcmeServer := r.apiClient.AcmeApi.AddAcmeServer(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAcmeServer = apiCreateAcmeServer.AcmeServer(*createAcmeServer)
	listenerResponse, httpResp, err := r.apiClient.AcmeApi.AddAcmeServerExecute(apiCreateAcmeServer)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the AcmeServer", err, httpResp)
		return
	}
	responseJson, err := listenerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state acmeserversResourceModel

	readAcmeServerResponse(ctx, listenerResponse, &state, &plan, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeserversResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAcmeServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAcmeServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state acmeserversResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAcmeServer, httpResp, err := apiClient.AcmeApi.GetAcmeServer(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an AcmeServer", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAcmeServer.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAcmeServerResponse(ctx, apiReadAcmeServer, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *acmeserversResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAcmeServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAcmeServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	tflog.Error(ctx, "This resource does not support updating. Please delete the resource and recreate to desired configuration.")
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *acmeserversResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAcmeServer(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAcmeServer(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state acmeserversResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := apiClient.AcmeApi.DeleteAcmeServer(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an AcmeServer", err, httpResp)
		return
	}

}

func (r *acmeserversResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
