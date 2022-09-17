package provider

import (
	"context"
	"fmt"
	"terraform-provider-msgraph/internal/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ApplicationWebDataSource{}

func NewApplicationWebDataSource() datasource.DataSource {
	return &ApplicationWebDataSource{}
}

// ApplicationWebDataSource defines the data source implementation.
type ApplicationWebDataSource struct {
	client *msgraph.Client
}

// ApplicationWebDataSourceModel describes the data source data model.
type ApplicationWebDataSourceModel struct {
	AppID       types.String `tfsdk:"app_id"`
	Id          types.String `tfsdk:"id"`
	Application types.Object `tfsdk:"application"`
}

func (d *ApplicationWebDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_web"
}

func (d *ApplicationWebDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Application data source",

		Attributes: map[string]tfsdk.Attribute{
			"app_id": {
				MarkdownDescription: "Application Client Id",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"application": {
				MarkdownDescription: "Application data",
				Type: types.ObjectType{AttrTypes: map[string]attr.Type{
					"appId":       types.StringType,
					"displayName": types.StringType,
					"id":          types.StringType,
					"web": types.ObjectType{AttrTypes: map[string]attr.Type{

						"implicitGrantSettings": types.ObjectType{AttrTypes: map[string]attr.Type{
							"enableAccessTokenIssuance": types.BoolType,
							"enableIdTokenIssuance":     types.BoolType,
						}},
						"redirectUris": types.ListType{
							ElemType: types.StringType,
						},
					}},
				}},
				Computed: true,
			},
		},
	}, nil
}

func (d *ApplicationWebDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*msgraph.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ApplicationWebDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationWebDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	d.client.GraphAccess()
	application := d.client.GetApplication(data.AppID.Value)

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: data.AppID.Value}

	redirectUris := make([]attr.Value, 0)
	for i := 0; i < len(application.Web.RedirectUris); i++ {
		redirectUris = append(redirectUris, types.String{Value: application.Web.RedirectUris[i]})
	}

	data.Application = types.Object{
		Unknown: false,
		Null:    false,
		Attrs: map[string]attr.Value{
			"appId":       types.String{Value: application.AppID},
			"displayName": types.String{Value: application.DisplayName},
			"id":          types.String{Value: application.ID},
			"web": types.Object{
				Unknown: false,
				Null:    false,
				Attrs: map[string]attr.Value{
					"implicitGrantSettings": types.Object{
						Attrs: map[string]attr.Value{
							"enableAccessTokenIssuance": types.Bool{Value: application.Web.ImplicitGrantSettings.EnableAccessTokenIssuance},
							"enableIdTokenIssuance":     types.Bool{Value: application.Web.ImplicitGrantSettings.EnableIDTokenIssuance},
						},
						AttrTypes: map[string]attr.Type{
							"enableAccessTokenIssuance": types.BoolType,
							"enableIdTokenIssuance":     types.BoolType,
						},
					},
					"redirectUris": types.List{
						Unknown:  false,
						Null:     false,
						Elems:    redirectUris,
						ElemType: types.StringType,
					},
				},
				AttrTypes: map[string]attr.Type{
					"implicitGrantSettings": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"enableAccessTokenIssuance": types.BoolType,
							"enableIdTokenIssuance":     types.BoolType,
						},
					},
					"redirectUris": types.ListType{
						ElemType: types.StringType,
					},
				},
			},
		},
		AttrTypes: data.Application.AttrTypes,
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
