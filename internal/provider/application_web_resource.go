package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-msgraph/internal/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ApplicationWebResource{}
var _ resource.ResourceWithImportState = &ApplicationWebResource{}

func NewApplicationWebResource() resource.Resource {
	return &ApplicationWebResource{}
}

// ApplicationWebResource defines the resource implementation.
type ApplicationWebResource struct {
	client *msgraph.Client
}

// ApplicationWebResourceModel describes the resource data model.
type ApplicationWebResourceModel struct {
	AppID       types.String `tfsdk:"app_id"`
	RedirectUri types.String `tfsdk:"redirect_uri"`
	Id          types.String `tfsdk:"id"`
	Application types.Object `tfsdk:"application"`
}

func (r *ApplicationWebResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_web"
}

func (r *ApplicationWebResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Application Web RedirectURI config resource",

		Attributes: map[string]tfsdk.Attribute{
			"app_id": {
				MarkdownDescription: "Application Client ID",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"redirect_uri": {
				MarkdownDescription: "Redirect URI attribute",
				Required:            true,
				Optional:            false,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
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

func (r *ApplicationWebResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*msgraph.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ApplicationWebResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ApplicationWebResourceModel

	tflog.Trace(ctx, "========= IN CREATE ==========")

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: data.AppID.Value}
	r.client.GraphAccess()
	application := r.client.GetApplication(data.AppID.Value)

	app, err := json.Marshal(application)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application data, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "in create json "+string(app))

	if r.client.CheckRedirectURI(*application, data.RedirectUri.Value) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("RedirectURI is present in Application: %s", err))
		return
	}

	err = r.client.PatchWebAddRedirectURI(*application, data.RedirectUri.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to patch application data, got error: %s", err))
		return
	}

	application = r.client.GetApplication(data.AppID.Value)

	redirectUris := make([]attr.Value, 0)
	for i := 0; i < len(application.Web.RedirectUris); i++ {
		redirectUris = append(redirectUris, types.String{Value: application.Web.RedirectUris[i]})
	}

	updateDataApplicationObject(data, application, redirectUris)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "========= IN CREATE END ==========")
}

func (r *ApplicationWebResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationWebResourceModel

	tflog.Trace(ctx, "========= IN READ ==========")

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.String{Value: data.AppID.Value}
	r.client.GraphAccess()
	application := r.client.GetApplication(data.AppID.Value)

	redirectUris := make([]attr.Value, 0)
	for i := 0; i < len(application.Web.RedirectUris); i++ {
		redirectUris = append(redirectUris, types.String{Value: application.Web.RedirectUris[i]})
	}

	updateDataApplicationObject(data, application, redirectUris)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "========= IN READ END ==========")
}

func (r *ApplicationWebResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ApplicationWebResourceModel

	tflog.Trace(ctx, "========= IN UPDATE ==========")

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *ApplicationWebResourceModel
	// Read Terraform State data into the model
	_ = req.State.Get(ctx, &state)

	r.client.GraphAccess()
	stateApplication := r.client.GetApplication(state.AppID.Value)

	if state.RedirectUri.Value != data.RedirectUri.Value {
		tflog.Trace(ctx, fmt.Sprintf("State Value %s == %s Plan Value", state.RedirectUri.Value, data.RedirectUri.Value))
		err := r.client.PatchWebRemoveRedirectURI(*stateApplication, state.RedirectUri.Value)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Update application data (delete uri), got error: %s", err))
			return
		}

		newApplication := r.client.GetApplication(state.AppID.Value)
		err = r.client.PatchWebAddRedirectURI(*newApplication, data.RedirectUri.Value)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Update application data (create uri), got error: %s", err))
			return
		}
	}

	application := r.client.GetApplication(state.AppID.Value)

	redirectUris := make([]attr.Value, 0)
	for i := 0; i < len(application.Web.RedirectUris); i++ {
		redirectUris = append(redirectUris, types.String{Value: application.Web.RedirectUris[i]})
	}
	updateDataApplicationObject(data, application, redirectUris)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "========= IN UPDATE END ==========")
}

func (r *ApplicationWebResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationWebResourceModel

	tflog.Trace(ctx, "========= IN DELETE ==========")
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.String{Value: data.AppID.Value}
	r.client.GraphAccess()
	application := r.client.GetApplication(data.AppID.Value)

	app, err := json.Marshal(application)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application data, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "in delete json "+string(app))
	if r.client.CheckRedirectURI(*application, data.RedirectUri.Value) {
		r.client.PatchWebRemoveRedirectURI(*application, data.RedirectUri.Value)
	}

	tflog.Trace(ctx, "========= IN DELETE END ==========")
}

func (r *ApplicationWebResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updateDataApplicationObject(data *ApplicationWebResourceModel, application *msgraph.Application, redirectUris []attr.Value) {
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
}
