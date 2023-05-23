package trustedCertificateGroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	_ resource.Resource                = &trustedCertificateGroupResource{}
	_ resource.ResourceWithConfigure   = &trustedCertificateGroupResource{}
	_ resource.ResourceWithImportState = &trustedCertificateGroupResource{}
)

// TrustedCertificateGroupResource is a helper function to simplify the provider implementation.
func TrustedCertificateGroupResource() resource.Resource {
	return &trustedCertificateGroupResource{}
}

// trustedCertificateGroupResource is the resource implementation.
type trustedCertificateGroupResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type trustedCertificateGroupResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	CertIds                    types.Set    `tfsdk:"cert_ids"`
	IgnoreAllCertificateErrors types.Bool   `tfsdk:"ignore_all_certificate_errors"`
	Name                       types.String `tfsdk:"name"`
	RevocationChecking         types.Object `tfsdk:"revocation_checking"`
	SkipCertificateDateCheck   types.Bool   `tfsdk:"skip_certificate_date_check"`
	SystemGroup                types.Bool   `tfsdk:"system_group"`
	UseJavaTrustStore          types.Bool   `tfsdk:"use_java_trust_store"`
}

// GetSchema defines the schema for the resource.
func (r *trustedCertificateGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	trustedCertificateGroupResourceSchema(ctx, req, resp, false)
}

func trustedCertificateGroupResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a TrustedCertificateGroup.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cert_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"ignore_all_certificate_errors": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_java_trust_store": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"revocation_checking": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"deny_revocation_status_unknown": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"crl_checking": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"ocsp": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"support_disordered_chain": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"skip_trust_anchors": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
				},
			},
			"skip_certificate_date_check": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"system_group": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name"})
	}
	resp.Schema = schema
}
func addOptionalTrustedCertificateGroupFields(ctx context.Context, addRequest *client.TrustedCertificateGroup, plan trustedCertificateGroupResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = internaltypes.StringToInt64Pointer(plan.Id)
	}

	if internaltypes.IsDefined(plan.UseJavaTrustStore) {
		addRequest.UseJavaTrustStore = plan.UseJavaTrustStore.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.SystemGroup) {
		addRequest.SystemGroup = plan.SystemGroup.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.IgnoreAllCertificateErrors) {
		addRequest.IgnoreAllCertificateErrors = plan.IgnoreAllCertificateErrors.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.CertIds) {
		var slice []int64
		plan.CertIds.ElementsAs(ctx, &slice, false)
		addRequest.CertIds = slice
	}

	if internaltypes.IsDefined(plan.SkipCertificateDateCheck) {
		addRequest.SkipCertificateDateCheck = plan.SkipCertificateDateCheck.ValueBoolPointer()
	}

	addRequest.RevocationChecking = client.NewRevocationChecking()
	if internaltypes.IsDefined(plan.RevocationChecking) {
		rC := plan.RevocationChecking.Attributes()
		addRequest.RevocationChecking.CrlChecking = internaltypes.InterfaceBoolPointerValue(internaltypes.ConvertToPrimitive(rC["crl_checking"]))
		addRequest.RevocationChecking.DenyRevocationStatusUnknown = internaltypes.InterfaceBoolPointerValue(internaltypes.ConvertToPrimitive(rC["deny_revocation_status_unknown"]))
		addRequest.RevocationChecking.Ocsp = internaltypes.InterfaceBoolPointerValue(internaltypes.ConvertToPrimitive(rC["ocsp"]))
		addRequest.RevocationChecking.SkipTrustAnchors = internaltypes.InterfaceBoolPointerValue(internaltypes.ConvertToPrimitive(rC["skip_trust_anchors"]))
		addRequest.RevocationChecking.SupportDisorderedChain = internaltypes.InterfaceBoolPointerValue(internaltypes.ConvertToPrimitive(rC["support_disordered_chain"]))
	}

	return nil
}

// Metadata returns the resource type name.
func (r *trustedCertificateGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_certificate_groups"
}

func (r *trustedCertificateGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readTrustedCertificateGroupResponse(ctx context.Context, r *client.TrustedCertificateGroup, state *trustedCertificateGroupResourceModel, expectedValues *trustedCertificateGroupResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(internaltypes.Int64PointerToString(*r.Id))
	state.CertIds = internaltypes.GetInt64SetOrNull(r.CertIds)
	state.IgnoreAllCertificateErrors = internaltypes.BoolTypeOrNil(r.IgnoreAllCertificateErrors)
	state.Name = types.StringValue(r.Name)
	state.SkipCertificateDateCheck = internaltypes.BoolTypeOrNil(r.SkipCertificateDateCheck)
	state.SystemGroup = internaltypes.BoolTypeOrNil(r.SystemGroup)
	state.UseJavaTrustStore = internaltypes.BoolTypeOrNil(r.UseJavaTrustStore)

	attrTypes := map[string]attr.Type{
		"crl_checking":                   basetypes.BoolType{},
		"ocsp":                           basetypes.BoolType{},
		"deny_revocation_status_unknown": basetypes.BoolType{},
		"support_disordered_chain":       basetypes.BoolType{},
		"skip_trust_anchors":             basetypes.BoolType{},
	}

	getRc := r.GetRevocationChecking()
	attrValues := map[string]attr.Value{
		"crl_checking":                   internaltypes.BoolTypeOrNil(getRc.CrlChecking),
		"ocsp":                           internaltypes.BoolTypeOrNil(getRc.Ocsp),
		"deny_revocation_status_unknown": internaltypes.BoolTypeOrNil(getRc.DenyRevocationStatusUnknown),
		"support_disordered_chain":       internaltypes.BoolTypeOrNil(getRc.SupportDisorderedChain),
		"skip_trust_anchors":             internaltypes.BoolTypeOrNil(getRc.SkipTrustAnchors),
	}

	revocationChecking := internaltypes.MaptoObjValue(attrTypes, attrValues, *diagnostics)
	state.RevocationChecking = revocationChecking
}

func (r *trustedCertificateGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan trustedCertificateGroupResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTrustedCertificateGroup := client.NewTrustedCertificateGroup(plan.Name.ValueString())
	err := addOptionalTrustedCertificateGroupFields(ctx, createTrustedCertificateGroup, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TrustedCertificateGroup", err.Error())
		return
	}
	requestJson, err := createTrustedCertificateGroup.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateTrustedCertificateGroup := r.apiClient.TrustedCertificateGroupsApi.AddTrustedCertificateGroup(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateTrustedCertificateGroup = apiCreateTrustedCertificateGroup.TrustedCertificateGroup(*createTrustedCertificateGroup)
	trustedCertificateGroupResponse, httpResp, err := r.apiClient.TrustedCertificateGroupsApi.AddTrustedCertificateGroupExecute(apiCreateTrustedCertificateGroup)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the TrustedCertificateGroup", err, httpResp)
		return
	}
	responseJson, err := trustedCertificateGroupResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state trustedCertificateGroupResourceModel

	readTrustedCertificateGroupResponse(ctx, trustedCertificateGroupResponse, &state, &plan, &diags)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trustedCertificateGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTrustedCertificateGroup(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readTrustedCertificateGroup(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state trustedCertificateGroupResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadTrustedCertificateGroup, httpResp, err := apiClient.TrustedCertificateGroupsApi.GetTrustedCertificateGroup(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a TrustedCertificateGroup", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadTrustedCertificateGroup.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readTrustedCertificateGroupResponse(ctx, apiReadTrustedCertificateGroup, &state, &state, &diags)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustedCertificateGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTrustedCertificateGroup(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTrustedCertificateGroup(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan trustedCertificateGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state trustedCertificateGroupResourceModel
	req.State.Get(ctx, &state)
	UpdateTrustedCertificateGroup := apiClient.TrustedCertificateGroupsApi.UpdateTrustedCertificateGroup(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	CreateUpdateRequest := client.NewTrustedCertificateGroup(plan.Name.ValueString())
	err := addOptionalTrustedCertificateGroupFields(ctx, CreateUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for TrustedCertificateGroup", err.Error())
		return
	}
	requestJson, err := CreateUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	UpdateTrustedCertificateGroup = UpdateTrustedCertificateGroup.TrustedCertificateGroup(*CreateUpdateRequest)
	updateTrustedCertificateGroupResponse, httpResp, err := apiClient.TrustedCertificateGroupsApi.UpdateTrustedCertificateGroupExecute(UpdateTrustedCertificateGroup)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating TrustedCertificateGroup", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateTrustedCertificateGroupResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readTrustedCertificateGroupResponse(ctx, updateTrustedCertificateGroupResponse, &state, &plan, &diags)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *trustedCertificateGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteTrustedCertificateGroup(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteTrustedCertificateGroup(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state trustedCertificateGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.TrustedCertificateGroupsApi.DeleteTrustedCertificateGroup(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a TrustedCertificateGroup", err, httpResp)
		return
	}

}

func (r *trustedCertificateGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
