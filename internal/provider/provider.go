package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &pingaccessProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &pingaccessProvider{}
}

// pingaccessProvider is the provider implementation.
type pingaccessProvider struct{}

// Metadata returns the provider type name.
func (p *pingaccessProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingaccess"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *pingaccessProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

// Configure prepares a PingAccess Provider API client for data sources and resources.
func (p *pingaccessProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *pingaccessProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *pingaccessProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
