package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	// "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	// "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// client "github.com/pingidentity/pingaccess-go-client"
)

// Ensure the implementation satisfies the expected interfacesß
var (
	_ provider.Provider = &pingaccessProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &pingaccessProvider{}
}

// PingAccess ProviderModel maps provider schema data to a Go type.
type pingaccessProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// pingaccessProvider is the provider implementation.
type pingaccessProvider struct{}

// Metadata returns the provider type name.
func (p *pingaccessProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingaccess"
}

// GetSchema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *pingaccessProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
                Optional: true,
            },
            "username": schema.StringAttribute{
                Optional: true,
            },
            "password": schema.StringAttribute{
                Optional:  true,
                Sensitive: true,
            },
        },
    }
}


func (p *pingaccessProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config pingaccessProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown PingAccess API Host",
			"The provider cannot create the PingAccess API client as there is an unknown configuration value for the PingAccess API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown PingAccess API Username",
			"The provider cannot create the PingAccess API client as there is an unknown configuration value for the PingAccess API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown PingAccess API Password",
			"The provider cannot create the PingAccess API client as there is an unknown configuration value for the PingAccess API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("PingAccess_HOST")
	username := os.Getenv("PingAccess_USERNAME")
	password := os.Getenv("PingAccess_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing PingAccess API Host",
			"The provider cannot create the PingAccess API client as there is a missing or empty value for the PingAccess API host. "+
				"Set the host value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing PingAccess API Username",
			"The provider cannot create the PingAccess API client as there is a missing or empty value for the PingAccess API username. "+
				"Set the username value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing PingAccess API Password",
			"The provider cannot create the PingAccess API client as there is a missing or empty value for the PingAccess API password. "+
				"Set the password value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new PingAccess client using the configuration values
	// client, err := provider.NewClient(&host, &username, &password)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to Create PingAccess API Client",
	// 		"An unexpected error occurred when creating the PingAccess API client. "+
	// 			"If the error is not clear, please contact the provider developers.\n\n"+
	// 			"PingAccess Client Error: "+err.Error(),
	// 	)
	// 	return
	// }

	// Make the PingAccess client available during DataSource and Resource
	// type Configure methods.
// 	resp.DataSourceData = client
// 	resp.ResourceData = client
 }

// DataSources defines the data sources implemented in the provider.
func (p *pingaccessProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *pingaccessProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEnginelistenerResource,
	}
}
