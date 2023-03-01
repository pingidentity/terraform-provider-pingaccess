package enginelistener


import (
    "context"
    "strconv"
    "time"

    "github.com/hashicorp-demoapp/hashicups-client-go"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"

)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource = &enginelistenerResource{}
    _ resource.ResourceWithConfigure = &enginelistenerResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewEnginelistenerResource() resource.Resource {
    return &enginelistenerResource{}
}

// orderResource is the resource implementation.
type enginelistenerResource struct{
    apiClient      *client.APIClient
    Items       []enginelistenerModel `tfsdk:"items"`
    LastUpdated types.String     `tfsdk:"last_updated"`
}
type enginelistenerItemModel struct {
    ID          types.Int64   `tfsdk:"id"`
    Name        types.String  `tfsdk:"name"`
    Port        types.String  `tfsdk:"port"`
    Secure      types.String  `tfsdk:"secure"`
    TrustedCertificateGroupId types.Int64 `tfsdk:"trustedCertificateGroupId`
}
// Metadata returns the resource type name.
func (r *enginelistenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_order"
}

// GetSchema defines the schema for the resource.
func (r *enginelistenerResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{}, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *enginelistenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *enginelistenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *enginelistenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *enginelistenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *enginelistenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.apiClient = req.ProviderData.(*client.APIClient)
}
// Schema defines the schema for the resource.
func (r *enginelistenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "items": schema.ListNestedAttribute{
                Required: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.Int64Attribute{
                            Required: true,
                        },
                        "name": schema.StringAttribute{
                            Computed: true,
                        },
                        "port": schema.int64Attribute{
                            Computed: true,
                        },
                        "secure": schema.StringAttribute{
                            Computed: true,
                        },
                        "trustedCertificateGroupId": schema.Int64Attribute{
                            Computed: true
                        },
                        "last_updated": schema.StringAttribute{
                            Computed: true,
                        },
                    },
                },
            },
        },
     },
 }

// Create a new resource
func (r *enginelistenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan enginelistenerResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    var items []client.Enginelistener
    for _, item := range plan.Items {
        items = append(items, client.EnginelistenerItem{
            ID: int(item.ID.ValueInt64()),
            Name: string(item.Name.ValueString()),
            Port: int(item.Port.ValueString()),
            Secure: string(item.Secure.ValueString()),
            TrustedCertificateGroupId: int(item.TrustedCertificateGroupId.ValueInt64()),
        })
    }

    // Create new order
    createlistener, err := r.apiclient.CreateOrder(items)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error while creating the listener",
            "Could not create order, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    for EnginelistenerItemIndex, EnginelistenerItem := range createlistener.Items {
        plan.Items[EnginelistenerItemIndex] = enginelistenerModel{
                id:                         types.Int64Value(EnginelistenerItem.ID),
                name:                       types.StringValue(orderItemEnginelistenerItem.Name),
                port:                       types.StringValue(orderItemEnginelistenerItem.Port),
                secure:                     types.StringValue(orderItemEnginelistenerItem.Secure),
                trustedCertificateGroupId:  types.Int64Value(EnginelistenerItem.TrustedCertificateGroupId),
        }
    }
  

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read resource information
func (r *enginelistenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
        var state enginelistenerResourceModel
        diags := req.State.Get(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
    
        // Get refreshed order value from HashiCups
        readlistener, err := r.apiclient.GetOrder(state.ID.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(
                "Error Reading EngineListener Order",
                "Could not read EngineListener order ID's "+state.ID.ValueString()+": "+err.Error(),
            )
            return
        }
    
        // Overwrite items with refreshed state
        state.Items = []enginelistenerModel {}
        for _, item := range readlistener.Items {
            state.Items = append(state.Items, enginelistenerModel{
                id:                         types.Int64Value(EnginelistenerItem.ID),
                name:                       types.StringValue(orderItemEnginelistenerItem.Name),
                port:                       types.StringValue(orderItemEnginelistenerItem.Port),
                secure:                     types.StringValue(orderItemEnginelistenerItem.Secure),
                trustedCertificateGroupId:  types.Int64Value(EnginelistenerItem.TrustedCertificateGroupId),
            })
        }
        plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
        // Set refreshed state
        diags = resp.State.Set(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
    }
    
func (r *enginelistenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan enginelistenerResource
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    var listenerItems []client.Enginelistener
    for _, item := range plan.Items {
        listenerItems = append(listenerItems, lient.Enginelistener{
            id:                         types.Int64Value(EnginelistenerItem.ID),
            name:                       types.StringValue(orderItemEnginelistenerItem.Name),
            port:                       types.StringValue(orderItemEnginelistenerItem.Port),
            secure:                     types.StringValue(orderItemEnginelistenerItem.Secure),
            trustedCertificateGroupId:  types.Int64Value(EnginelistenerItem.TrustedCertificateGroupId),
        })
    }

    // Update existing order
    _, err := r.apiclient.UpdateOrder(plan.ID.ValueString(), listenerItems)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Updating enginelistener",
            "Could not update enginelistener, unexpected error: "+err.Error(),
        )
        return
    }

    // Fetch updated items from GetOrder as UpdateOrder items are not
    // populated.
    updatelistener, err := r.apiclient.GetOrder(plan.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading enginelistener",
            "Could not read enginelistener ID "+plan.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Update resource state with updated items and timestamp
    plan.Items = []enginelistenerModel{}
    for _, item := range updatelistener.Items {
        plan.Items = append(plan.Items, enginelistenerModel{
            id:                         types.Int64Value(EnginelistenerItem.ID),
            name:                       types.StringValue(orderItemEnginelistenerItem.Name),
            port:                       types.StringValue(orderItemEnginelistenerItem.Port),
            secure:                     types.StringValue(orderItemEnginelistenerItem.Secure),
            trustedCertificateGroupId:  types.Int64Value(EnginelistenerItem.TrustedCertificateGroupId),
        })
    }
    plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}
func (r *enginelistenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
        var enginelistenerResourceModel
        diags := req.State.Get(ctx, &state)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
    
        // Delete existing order
        err := r.apiclient.DeleteOrder(state.ID.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(
                "Error Deleting enginelistener",
                "Could not delete enginelistener, unexpected error: "+err.Error(),
            )
            return
        }
    }
    