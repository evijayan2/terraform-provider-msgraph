package provider

import (
	"context"
	"os"
	"terraform-provider-msgraph/internal/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure MsgraphProvider satisfies various provider interfaces.
var _ provider.Provider = &MsgraphProvider{}
var _ provider.ProviderWithMetadata = &MsgraphProvider{}

// MsgraphProvider defines the provider implementation.
type MsgraphProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MsgraphProviderModel describes the provider data model.
type MsgraphProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	TenantID     types.String `tfsdk:"tenant_id"`
	Scope        types.String `tfsdk:"scope"`
	GrantType    types.String `tfsdk:"grant_type"`
	AuthHost     types.String `tfsdk:"auth_host"`
	GraphHost    types.String `tfsdk:"graph_host"`
}

func (p *MsgraphProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msgraph"
	resp.Version = p.version
}

func (p *MsgraphProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"client_id": {
				MarkdownDescription: "Client id",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"client_secret": {
				MarkdownDescription: "Client secret",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"tenant_id": {
				MarkdownDescription: "Tenant Id",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"scope": {
				MarkdownDescription: "Scope",
				Optional:            true,
				Type:                types.StringType,
			},
			"grant_type": {
				MarkdownDescription: "Grant Type",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"auth_host": {
				MarkdownDescription: "MS Auth host",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"graph_host": {
				MarkdownDescription: "MS Graph host",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *MsgraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data MsgraphProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.ClientID.IsNull() {
		data.ClientID.Value = os.Getenv("CLIENT_ID")
	}

	if data.ClientSecret.IsNull() {
		resp.Diagnostics.AddError("Client secret is missing", "this is critical")
	}

	if data.Scope.IsNull() {
		data.Scope.Value = "https://graph.microsoft.com/.default"
	}

	config := msgraph.ClientConfiguration{
		ClientID:     data.ClientID.Value,
		ClientSecret: data.ClientSecret.Value,
		TenantID:     data.TenantID.Value,
		Scope:        data.Scope.Value,
		GrantType:    data.GrantType.Value,
		AuthHost:     data.AuthHost.Value,
		GraphHost:    data.GraphHost.Value,
		UserAgent:    GetUserAgent(),
	}

	// Example client configuration for data sources and resources
	client := msgraph.NewClient(config)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MsgraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationWebResource,
	}
}

func (p *MsgraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationWebDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MsgraphProvider{
			version: version,
		}
	}
}

func GetUserAgent() string {
	return "terraform-provider-msgraph"
}
