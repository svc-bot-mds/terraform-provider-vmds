package mds

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"time"
)

const (
	defaultUserCreateTimeout = 2 * time.Minute
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *mds.Client
}

type userResourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Email        types.String   `tfsdk:"email"`
	Status       types.String   `tfsdk:"status"`
	Username     types.String   `tfsdk:"username"`
	PolicyIds    types.Set      `tfsdk:"policy_ids"`
	RoleIds      []string       `tfsdk:"role_ids"`
	ServiceRoles types.List     `tfsdk:"service_roles"`
	OrgRoles     types.List     `tfsdk:"org_roles"`
	Tags         types.Set      `tfsdk:"tags"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

type RolesModel struct {
	ID   types.String `tfsdk:"role_id"`
	Name types.String `tfsdk:"name"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *userResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents an User registered on MDS, can be used to create/update/delete/import an user.\n" +
			"## Notes\n" +
			fmt.Sprintf("- Default timeout for creation is `%v`.", defaultUserCreateTimeout),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID after creating an user, and can be passed to import an existing user from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "Updating the email results in deletion of existing user and new user with updated email/name is created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description: "Active status of user on MDS.",
				Computed:    true,
				Default:     stringdefault.StaticString("INVITED"),
			},
			"username": schema.StringAttribute{
				Description: "Short name of user.",
				Computed:    true,
			},
			"policy_ids": schema.SetAttribute{
				Description: "IDs of service policies to be associated with user.",
				Optional:    true,
				Computed:    false,
				ElementType: types.StringType,
			},
			"role_ids": schema.SetAttribute{
				MarkdownDescription: "One or more of (Admin, Developer, Viewer, Operator, Compliance Manager). Please make use of `datasource_roles` to get role_ids.",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				Description: "Tags or labels to categorise users for ease of finding.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"service_roles": schema.ListNestedAttribute{
				Description: "Roles that determines access level inside services on MDS.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							Description: "ID of the role.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the role.",
							Computed:    true,
						},
					},
				},
			},
			"org_roles": schema.ListNestedAttribute{
				Description: "Roles that determines access level of the user on MDS.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							Description: "ID of the role.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the role.",
							Computed:    true,
						},
					},
				},
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultUserCreateTimeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()
	rolesReq := make([]customer_metadata.RolesRequest, len(plan.RoleIds))
	for i, roleId := range plan.RoleIds {
		rolesReq[i] = customer_metadata.RolesRequest{
			RoleId: roleId,
		}
	}

	// Generate API request body from plan
	userRequest := customer_metadata.MdsCreateUserRequest{
		Usernames:    []string{plan.Email.ValueString()},
		ServiceRoles: rolesReq,
	}
	plan.Tags.ElementsAs(ctx, &userRequest.Tags, true)
	plan.PolicyIds.ElementsAs(ctx, &userRequest.PolicyIds, true)
	if err := r.client.CustomerMetadata.CreateMdsUser(&userRequest); err != nil {
		resp.Diagnostics.AddError(
			"Submitting request to create User",
			"Could not create User, unexpected error: "+err.Error(),
		)
		return
	}

	users, err := r.client.CustomerMetadata.GetMdsUsers(&customer_metadata.MdsUsersQuery{
		Emails: []string{plan.Email.ValueString()},
	})

	if err != nil {
		resp.Diagnostics.AddError("Fetching user",
			"Could not fetch users, unexpected error: "+err.Error(),
		)
		return
	}
	if users.Page.TotalElements == 0 {
		resp.Diagnostics.AddError("Fetching user",
			fmt.Sprintf("Could not find any user by email [%s], server error must have occurred while creating user.", plan.Email.ValueString()),
		)
		return
	}

	if len(*users.Get()) <= 0 {
		resp.Diagnostics.AddError("Fetching User",
			"Unable to fetch the created user",
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	createdUser := &(*users.Get())[0]

	if saveFromUserResponse(&ctx, &resp.Diagnostics, &plan, createdUser) != 0 {
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

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRequest := customer_metadata.MdsUserUpdateRequest{}
	plan.Tags.ElementsAs(ctx, &updateRequest.Tags, true)
	plan.PolicyIds.ElementsAs(ctx, &updateRequest.PolicyIds, true)
	if plan.Status.ValueString() != "INVITED" {
		rolesReq := make([]customer_metadata.RolesRequest, len(plan.RoleIds))
		for i, roleId := range plan.RoleIds {
			rolesReq[i] = customer_metadata.RolesRequest{
				RoleId: roleId,
			}
		}
		tflog.Info(ctx, "Setting serviceRoles for update: ", map[string]interface{}{
			"roles": rolesReq,
		})
		if len(rolesReq) > 0 {
			updateRequest.ServiceRoles = &rolesReq
		}
	}

	// Update existing user
	if err := r.client.CustomerMetadata.UpdateMdsUser(plan.ID.ValueString(), &updateRequest); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS User",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	user, err := r.client.CustomerMetadata.GetMdsUser(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching User",
			"Could not fetch users while updating, unexpected error: "+err.Error(),
		)
		return
	}

	//Update resource state with updated items and timestamp
	if saveFromUserResponse(&ctx, &resp.Diagnostics, &plan, user) != 0 {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *userResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state userResourceModel
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
	err := r.client.CustomerMetadata.DeleteMdsUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting MDS User",
			"Could not delete MDS User by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed cluster value from MDS
	user, err := r.client.CustomerMetadata.GetMdsUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS user",
			"Could not read MDS user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	if saveFromUserResponse(&ctx, &resp.Diagnostics, &state, user) != 0 {
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

func saveFromUserResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *userResourceModel, user *model.MdsUser) int8 {
	tflog.Info(*ctx, "Saving response to resourceModel state/plan", map[string]interface{}{"user": *user})

	roles, diags := convertFromRolesDto(ctx, &user.ServiceRoles)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.ServiceRoles = roles

	roles, diags = convertFromRolesDto(ctx, &user.OrgRoles)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.OrgRoles = roles
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}

	state.ID = types.StringValue(user.Id)
	state.Email = types.StringValue(user.Email)
	state.Status = types.StringValue(user.Status)
	tags, diags := types.SetValueFrom(*ctx, types.StringType, user.Tags)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.Tags = tags
	state.Username = types.StringValue(user.Name)

	return 0
}

func convertFromRolesDto(ctx *context.Context, roles *[]model.MdsRoleMini) (types.List, diag.Diagnostics) {
	tfRoleModels := make([]RolesModel, len(*roles))
	for i, role := range *roles {
		tfRoleModels[i] = RolesModel{
			ID:   types.StringValue(role.RoleID),
			Name: types.StringValue(role.Name),
		}
	}
	return types.ListValueFrom(*ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":    types.StringType,
		"role_id": types.StringType,
	}}, tfRoleModels)
}
