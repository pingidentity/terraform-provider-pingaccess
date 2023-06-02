package certificates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource              = &certificatesResource{}
	_ resource.ResourceWithConfigure = &certificatesResource{}
)

// CertificateResource is a helper function to simplify the provider implementation.
func CertificateResource() resource.Resource {
	return &certificatesResource{}
}

// certificatesResource is the resource implementation.
type certificatesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type certificatesResourceModel struct {
	Id       types.Int64  `tfsdk:"id"`
	Alias    types.String `tfsdk:"alias"`
	FileData types.String `tfsdk:"file_data"`
}

// GetSchema defines the schema for the resource.
func (r *certificatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificatesResourceSchema(ctx, req, resp, false)
}

func certificatesResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages Cetrificate Import.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"alias": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_data": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"alias", "file_data"})
	}
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *certificatesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

func (r *certificatesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readCertificateResponse(ctx context.Context, r *client.TrustedCert, state *certificatesResourceModel, expectedValues *certificatesResourceModel, diagnostics *diag.Diagnostics, createPlan types.String) {
	X509FileData := createPlan
	state.Id = internaltypes.Int64InterfaceTypeOrNil(*r.Id)
	state.Alias = types.StringValue(r.Alias)
	state.FileData = types.StringValue(X509FileData.ValueString())
}

func (r *certificatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan certificatesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createCertificate := client.NewX509FileImportDoc(plan.Alias.ValueString(), plan.FileData.ValueString())
	requestJson, err := createCertificate.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiCreateCertificate := r.apiClient.CertificatesApi.ImportTrustedCert(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateCertificate = apiCreateCertificate.X509File(*createCertificate)
	certificateResponse, httpResp, err := r.apiClient.CertificatesApi.ImportTrustedCertExecute(apiCreateCertificate)
	if httpResp.StatusCode != 200 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating a Certificate", err, httpResp)
		return
	}
	responseJson, err := certificateResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificatesResourceModel

	readCertificateResponse(ctx, certificateResponse, &state, &plan, &resp.Diagnostics, plan.FileData)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificate(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCertificate(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state certificatesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadCertificate, httpResp, err := apiClient.CertificatesApi.GetTrustedCert(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id)).Execute()

	if httpResp.StatusCode != 200 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a Certificate", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadCertificate.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCertificateResponse(ctx, apiReadCertificate, &state, &state, &resp.Diagnostics, state.FileData)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *certificatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCertificate(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCertificate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan certificatesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state certificatesResourceModel
	req.State.Get(ctx, &state)
	updateCertificate := apiClient.CertificatesApi.UpdateTrustedCert(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id))
	CreateUpdateRequest := client.NewX509FileImportDoc(plan.Alias.ValueString(), plan.FileData.ValueString())
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateCertificate = updateCertificate.X509File(*CreateUpdateRequest)
	updateCertificateResponse, httpResp, err := apiClient.CertificatesApi.UpdateTrustedCertExecute(updateCertificate)
	if httpResp.StatusCode != 200 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating a certificate", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateCertificateResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readCertificateResponse(ctx, updateCertificateResponse, &state, &plan, &resp.Diagnostics, plan.FileData)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *certificatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteCertificate(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteCertificate(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state certificatesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := apiClient.CertificatesApi.DeleteTrustedCert(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id)).Execute()
	if httpResp.StatusCode != 200 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Certificate", err, httpResp)
		return
	}

}
