package mds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cloudAccountResource{}
	_ resource.ResourceWithConfigure   = &cloudAccountResource{}
	_ resource.ResourceWithImportState = &cloudAccountResource{}
)

func NewCloudAccountResource() resource.Resource {
	return &cloudAccountResource{}
}

type cloudAccountResource struct {
	client *mds.Client
}

type CloudAccountResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProviderType   types.String `tfsdk:"provider_type"`
	Shared         types.Bool   `tfsdk:"shared"`
	Credential     types.String `tfsdk:"credentials"`
	Tags           types.Set    `tfsdk:"tags"`
	UserEmail      types.String `tfsdk:"user_email"`
	OrgId          types.String `tfsdk:"org_id"`
	CreatedBy      types.String `tfsdk:"created_by"`
	DataPlaneCount types.Int64  `tfsdk:"data_plane_count"`
}

func (r *cloudAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_account"
}

func (r *cloudAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *cloudAccountResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents a cloud account created on MDS, can be used to create/update/delete/import a cloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID after creating an cloud account, and can be passed to import an existing user from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name is readonly field while updating the certificate. Updating it will result in creating new account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_type": schema.StringAttribute{
				Description: "Provider Type of cloud account on MDS. Change not allowed after creation.",
				Required:    true,
			},
			"tags": schema.SetAttribute{
				Description: "Tags to set on this account.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_email": schema.StringAttribute{
				Description: "Email of the MDS User",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"shared": schema.BoolAttribute{
				Description: "Whether this account will be shared between multiple Organisations or not. Change not allowed after creation.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Org ID. Required for `shared` cloud account. Change not allowed after creation.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Holds the credentials associated with the cloud account.",
				Required:            true,
				Sensitive:           true,
			},
			"created_by": schema.StringAttribute{
				Description: "User which created this account.",
				Computed:    true,
			},
			"data_plane_count": schema.Int64Attribute{
				Description: "Total data planes associated with this account.",
				Computed:    true,
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *cloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan CloudAccountResourceModel

	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var cred infra_connector.CredentialModel
	if err := json.Unmarshal([]byte(plan.Credential.ValueString()), &cred); err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling Credential JSON",
			" Unexpected error: "+err.Error(),
		)
	} else {
		tflog.Info(ctx, "Successfully unmarshalled Credential JSON:", map[string]interface{}{"cred": cred})
	}

	// Generate API request body from plan
	cloudAccountRequest := &infra_connector.CloudAccountCreateRequest{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
		Credentials:  cred,
		Shared:       plan.Shared.ValueBool(),
	}
	plan.Tags.ElementsAs(ctx, &cloudAccountRequest.Tags, true)

	tflog.Info(ctx, "req param", map[string]interface{}{"create-request": cloudAccountRequest})
	cloudAccount, err := r.client.InfraConnector.CreateCloudAccount(cloudAccountRequest)
	if err != nil {
		apiErr := core.ApiError{}
		errors.As(err, &apiErr)
		resp.Diagnostics.AddError(
			"Submitting request to create cloud account",
			"There was some issue while creating the cloud account."+
				" Unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	if saveFromCloudAccountCreateResponse(&ctx, &resp.Diagnostics, &plan, cloudAccount) != 0 {
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

func (r *cloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var plan, state CloudAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if msg := r.validateUpdateRequest(&state, &plan); msg != "OK" {
		resp.Diagnostics.AddError(
			"Invalid Update", msg,
		)
	}

	request := infra_connector.CloudAccountUpdateRequest{}
	plan.Tags.ElementsAs(ctx, &request.Tags, true)
	var cred infra_connector.CredentialModel
	if err := json.Unmarshal([]byte(plan.Credential.ValueString()), &cred); err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling Credential JSON",
			" Unexpected error: "+err.Error(),
		)
	} else {
		tflog.Info(ctx, "Successfully unmarshalled Credential JSON:", map[string]interface{}{"cred": cred})
		request.Credentials = cred
	}

	// Update existing cloud account
	if err := r.client.InfraConnector.UpdateCloudAccount(plan.ID.ValueString(), &request); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS cloud account",
			"Could not update cloud account, unexpected error: "+err.Error(),
		)
		return
	}

	cloudAccount, err := r.client.InfraConnector.GetCloudAccount(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching cloud account",
			"Could not fetch cloud account while updating, unexpected error: "+err.Error(),
		)
		return
	}

	if saveFromCloudAccountCreateResponse(&ctx, &resp.Diagnostics, &plan, cloudAccount) != 0 {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *cloudAccountResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state CloudAccountResourceModel
	diags := request.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Submit request to delete MDS cloud Account
	err := r.client.InfraConnector.DeleteCloudAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting MDS cloud account",
			"Could not delete MDS cloud account by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *cloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *cloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state CloudAccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed cloud account value from MDS
	cloudAcct, err := r.client.InfraConnector.GetCloudAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS cloud account",
			"Could not read MDS cloud account ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if saveFromCloudAccountCreateResponse(&ctx, &resp.Diagnostics, &state, cloudAcct) != 0 {
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

func (r *cloudAccountResource) validateUpdateRequest(state *CloudAccountResourceModel, plan *CloudAccountResourceModel) string {
	if state.ProviderType != plan.ProviderType {
		return `Updating "provider_type" is not allowed`
	}
	if state.Shared != plan.Shared {
		return `Updating "shared" is not allowed`
	}
	if state.OrgId != plan.OrgId {
		return `Updating "org_id" is not allowed`
	}
	return "OK"
}

func saveFromCloudAccountCreateResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *CloudAccountResourceModel,
	response *model.MdsCloudAccount) int8 {
	state.ID = types.StringValue(response.Id)
	state.Name = types.StringValue(response.Name)
	state.UserEmail = types.StringValue(response.Email)
	state.ProviderType = types.StringValue(response.AccountType)
	state.OrgId = types.StringValue(response.OrgId)
	state.Shared = types.BoolValue(response.Shared)
	state.DataPlaneCount = types.Int64Value(response.DataPlaneCount)
	state.CreatedBy = types.StringValue(response.CreatedBy)
	list, diags := types.SetValueFrom(*ctx, types.StringType, response.Tags)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.Tags = list
	return 0
}
