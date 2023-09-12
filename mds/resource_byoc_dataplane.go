package mds

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &byocDataPlaneResource{}
	_ resource.ResourceWithConfigure   = &byocDataPlaneResource{}
	_ resource.ResourceWithImportState = &byocDataPlaneResource{}
)

func NewByocDataPlaneResourceResource() resource.Resource {
	return &byocDataPlaneResource{}
}

type byocDataPlaneResource struct {
	client *mds.Client
}

type byocDataPlaneResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AccountId            types.String `tfsdk:"account_id"`
	CertificateId        types.String `tfsdk:"certificate_id"`
	NodePoolType         types.String `tfsdk:"nodepool_type"`
	Region               types.String `tfsdk:"region"`
	Status               types.String `tfsdk:"status"`
	Version              types.String `tfsdk:"version"`
	Provider             types.String `tfsdk:"provider_name"`
	Certificate          types.Object `tfsdk:"certificate"`
	DataPlaneReleaseName types.String `tfsdk:"data_plane_release_name"`
}

type CertificateModel struct {
	Name       types.String `tfsdk:"name"`
	DomainName types.String `tfsdk:"domain_name"`
}

func (r *byocDataPlaneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_byoc_dataplane"
}

func (r *byocDataPlaneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mds.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *mds.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *byocDataPlaneResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")
	resp.Schema = schema.Schema{
		Description: "Represents a Dataplane on BYOC. Supported actions are Add and Delete ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID of the dataplane after creation, and can be used to import it from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Id of the selected Cloud Account",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the DataPlane",
				Required:    true,
			},
			"certificate_id": schema.StringAttribute{
				Description: "Id of the selected Certificate",
				Required:    true,
			},
			"nodepool_type": schema.StringAttribute{
				Description: "Selected T-shirt Size; Values can be 'regular' or 'large'",
				Required:    true,
			},
			"region": schema.StringAttribute{
				Description: "Selected Region",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the dataplane",
				Computed:    true,
			},
			"provider_name": schema.StringAttribute{
				Description: "Provider name",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "K8S version",
				Computed:    true,
			},
			"data_plane_release_name": schema.StringAttribute{
				Description: "Helm Release Name",
				Computed:    true,
			},
			"certificate": schema.SingleNestedAttribute{
				MarkdownDescription: "Certificate Details",
				Computed:            true,
				CustomType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"domain_name": types.StringType,
					},
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the Certificate",
						Computed:            true,
					},
					"domain_name": schema.StringAttribute{
						MarkdownDescription: "Domain Name of the Certificate",
						Computed:            true,
					},
				},
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *byocDataPlaneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan byocDataPlaneResourceModel
	diags := req.Plan.Get(ctx, &plan)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dataplaneRequest := infra_connector.ByocDataPlaneCreateRequest{
		Name:          plan.Name.ValueString(),
		CertificateId: plan.CertificateId.ValueString(),
		AccountId:     plan.AccountId.ValueString(),
		Region:        plan.Region.ValueString(),
		TshirtSize:    plan.NodePoolType.ValueString(),
	}

	tflog.Debug(ctx, "Create dataplane DTO", map[string]interface{}{"dto": dataplaneRequest})
	if _, err := r.client.InfraConnector.CreateDataPlane(&dataplaneRequest); err != nil {

		resp.Diagnostics.AddError(
			"Submitting request to create dataplane",
			"Could not create dataplane, unexpected error: "+err.Error(),
		)
		return
	}

	dataplanes, err := r.client.InfraConnector.GetDataPlaneByName(&infra_connector.ByocDataPlaneQuery{
		Name: dataplaneRequest.Name,
	})
	if err != nil {
		resp.Diagnostics.AddError("Fetching DataPlane",
			"Could not fetch data plane, unexpected error: "+err.Error(),
		)
		return
	}

	if len(*dataplanes.Get()) <= 0 {
		resp.Diagnostics.AddError("Fetching dataplane",
			"Unable to fetch the created dataplane",
		)
		return
	}
	createdDataPlane := &(*dataplanes.Get())[0]
	tflog.Debug(ctx, "Created dataplane DTO", map[string]interface{}{"dto": createdDataPlane})
	if saveFromDataPlaneResponse(&ctx, &resp.Diagnostics, &plan, createdDataPlane) != 0 {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Create")
}

func (r *byocDataPlaneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *byocDataPlaneResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state byocDataPlaneResourceModel
	diags := request.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Submit request to delete Byoc DataPlane
	err := r.client.InfraConnector.DeleteByocDataPlane(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting Byoc DataPlane",
			"Could not delete dataplane "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *byocDataPlaneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *byocDataPlaneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state byocDataPlaneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed dataplane value
	dataplane, err := r.client.InfraConnector.GetDataPlaneById(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading Byoc Dataplane",
			"Could not read dataplane "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	if saveFromDataPlaneResponse(&ctx, &resp.Diagnostics, &state, &dataplane) != 0 {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Read")
}

func saveFromDataPlaneResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *byocDataPlaneResourceModel, byocDataPlane *model.ByocDataPlane) int8 {
	tflog.Info(*ctx, "Saving response to resourceModel state/plan", map[string]interface{}{"byocDataPlane": *byocDataPlane})

	state.ID = types.StringValue(byocDataPlane.Id)
	state.Name = types.StringValue(byocDataPlane.Name)
	state.Status = types.StringValue(byocDataPlane.Status)
	state.Version = types.StringValue(byocDataPlane.K8SVersion)
	state.Provider = types.StringValue(byocDataPlane.Provider)
	state.DataPlaneReleaseName = types.StringValue(byocDataPlane.DataPlaneReleaseName)
	state.Status = types.StringValue(byocDataPlane.Status)
	state.Version = types.StringValue(byocDataPlane.K8SVersion)
	state.NodePoolType = types.StringValue(byocDataPlane.TshirtSize)

	certificateResponseModel := CertificateModel{
		Name:       types.StringValue(byocDataPlane.Certificate.Name),
		DomainName: types.StringValue(byocDataPlane.Certificate.DomainName),
	}

	certificateObject, diags := types.ObjectValueFrom(*ctx, state.Certificate.AttributeTypes(*ctx), certificateResponseModel)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}

	state.Certificate = certificateObject
	state.DataPlaneReleaseName = types.StringValue(byocDataPlane.DataPlaneReleaseName)
	return 0
}
