package mds

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/account_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"time"
)

const (
	defaultServiceAccCreateTimeout = 2 * time.Minute
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serviceAccountResource{}
	_ resource.ResourceWithConfigure   = &serviceAccountResource{}
	_ resource.ResourceWithImportState = &serviceAccountResource{}
)

func NewServiceAccountResource() resource.Resource {
	return &serviceAccountResource{}
}

type serviceAccountResource struct {
	client *mds.Client
}

type serviceAccountResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Status      types.String   `tfsdk:"status"`
	PolicyIds   types.Set      `tfsdk:"policy_ids"`
	Tags        types.Set      `tfsdk:"tags"`
	AccountType types.String   `tfsdk:"account_type"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

func (r *serviceAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account"
}

func (r *serviceAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *serviceAccountResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_type": schema.StringAttribute{
				Computed: true,
				Default:  stringdefault.StaticString(account_type.SERVICE_ACCOUNT),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Updating the name results in deletion of existing service account and new service account with updated name is created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"policy_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    false,
				ElementType: types.StringType,
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *serviceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan serviceAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultServiceAccCreateTimeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()
	// Generate API request body from plan
	svcAccountRequest := customer_metadata.MdsCreateSvcAccountRequest{
		Usernames:   []string{plan.Name.ValueString()},
		AccountType: account_type.SERVICE_ACCOUNT,
	}
	plan.Tags.ElementsAs(ctx, &svcAccountRequest.Tags, true)
	plan.PolicyIds.ElementsAs(ctx, &svcAccountRequest.PolicyIds, true)

	if err := r.client.CustomerMetadata.CreateMdsServiceAccount(&svcAccountRequest); err != nil {
		resp.Diagnostics.AddError(
			"Submitting request to create service account",
			"Could not create service account, unexpected error: "+err.Error(),
		)
		return
	}

	svcAccounts, err := r.client.CustomerMetadata.GetMdsServiceAccounts(&customer_metadata.MdsServiceAccountsQuery{
		AccountType: account_type.SERVICE_ACCOUNT,
		Name:        []string{plan.Name.ValueString()},
	})

	if err != nil {
		resp.Diagnostics.AddError("Fetching service account",
			"Could not fetch service accounts, unexpected error: "+err.Error(),
		)
		return
	}
	if svcAccounts.Page.TotalElements == 0 {
		resp.Diagnostics.AddError("Fetching svcAccounts",
			fmt.Sprintf("Could not find any svcAccounts by nasme [%s], server error must have occurred while creating svc account.", plan.Name.ValueString()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	createdSvcAccounts := &(*svcAccounts.Get())[0]

	if saveFromSvcAccountResponse(&ctx, &resp.Diagnostics, &plan, createdSvcAccounts) != 0 {
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

func (r *serviceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var plan serviceAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRequest := customer_metadata.MdsSvcAccountUpdateRequest{}
	plan.Tags.ElementsAs(ctx, &updateRequest.Tags, true)
	plan.PolicyIds.ElementsAs(ctx, &updateRequest.PolicyIds, true)

	// Update existing svc account
	if err := r.client.CustomerMetadata.UpdateMdsServiceAccount(plan.ID.ValueString(), &updateRequest); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS service account",
			"Could not update service account, unexpected error: "+err.Error(),
		)
		return
	}

	svcAccount, err := r.client.CustomerMetadata.GetMdsServiceAccount(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching svc account",
			"Could not fetch svc account while updating, unexpected error: "+err.Error(),
		)
		return
	}

	//Update resource state with updated items and timestamp
	if saveFromSvcAccountResponse(&ctx, &resp.Diagnostics, &plan, svcAccount) != 0 {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *serviceAccountResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state serviceAccountResourceModel
	diags := request.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	// Submit request to delete MDS Cluster
	err := r.client.CustomerMetadata.DeleteMdsServiceAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting MDS svc account",
			"Could not delete MDS svc account by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *serviceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *serviceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state serviceAccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed cluster value from MDS
	svcAcct, err := r.client.CustomerMetadata.GetMdsServiceAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS service account",
			"Could not read MDS service account ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	if saveFromSvcAccountResponse(&ctx, &resp.Diagnostics, &state, svcAcct) != 0 {
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

func saveFromSvcAccountResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *serviceAccountResourceModel, service_account *model.MdsServiceAccount) int8 {
	tflog.Info(*ctx, "Saving response to resourceModel state/plan", map[string]interface{}{"service_accounts": *service_account})

	state.ID = types.StringValue(service_account.Id)
	state.Name = types.StringValue(service_account.Name)
	state.Status = types.StringValue(service_account.Status)
	tags, diags := types.SetValueFrom(*ctx, types.StringType, service_account.Tags)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.Tags = tags
	state.AccountType = types.StringValue(account_type.SERVICE_ACCOUNT)

	return 0
}
