package mds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ProviderType types.String `tfsdk:"provider_type"`
	UserEmail    types.String `tfsdk:"user_email"`
	Credential   types.String `tfsdk:"credential"`
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
				Description: "Updating the name results in deletion of existing cloud account and new cloud account with updated name is created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_type": schema.StringAttribute{
				Description: "Provider Type of cloud account on MDS.",
				Required:    true,
			},
			"user_email": schema.StringAttribute{
				Description: "Email of the MDS User",
				Computed:    true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"credential": schema.StringAttribute{
				MarkdownDescription: "Holds the credentials associated with the cloud account.",
				Required:            true,
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
	}

	tflog.Info(ctx, "req param", map[string]interface{}{"reeed": cloudAccountRequest})
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
	if saveFromCloudAccountCreateResponse(&plan, cloudAccount) != 0 {
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
	var state CloudAccountResourceModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cred infra_connector.CredentialModel

	if err := json.Unmarshal([]byte(state.Credential.ValueString()), &cred); err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling Credential JSON",
			" Unexpected error: "+err.Error(),
		)
	} else {
		tflog.Info(ctx, "Successfully unmarshalled Credential JSON:", map[string]interface{}{"cred": cred})
	}

	// Update existing cloud account
	if err := r.client.InfraConnector.UpdateCloudAccount(state.ID.ValueString(), &cred); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS cloud account",
			"Could not update cloud account, unexpected error: "+err.Error(),
		)
		return
	}

	cloudAccount, err := r.client.InfraConnector.GetCloudAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching cloud account",
			"Could not fetch cloud account while updating, unexpected error: "+err.Error(),
		)
		return
	}

	if saveFromCloudAccountCreateResponse(&state, cloudAccount) != 0 {
		return
	}

	diags = resp.State.Set(ctx, state)
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

	if saveFromCloudAccountCreateResponse(&state, cloudAcct) != 0 {
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

func saveFromCloudAccountCreateResponse(state *CloudAccountResourceModel,
	cloudAccountCreateResponse *model.MdsCloudAccount) int8 {
	state.Name = types.StringValue(cloudAccountCreateResponse.Name)
	state.UserEmail = types.StringValue(cloudAccountCreateResponse.Email)
	state.ID = types.StringValue(cloudAccountCreateResponse.Id)
	state.ProviderType = types.StringValue(cloudAccountCreateResponse.AccountType)
	return 0
}
