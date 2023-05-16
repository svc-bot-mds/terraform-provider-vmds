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
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/policy_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/service_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

type policyResource struct {
	client *mds.Client
}

type policyResourceModel struct {
	ID              types.String           `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	ServiceType     types.String           `tfsdk:"service_type"`
	PermissionSpecs []*PermissionSpecModel `tfsdk:"permission_specs"`
	NetworkSpec     *NetworkSpecModel      `tfsdk:"network_spec"`
	ResourceIds     types.Set              `tfsdk:"resource_ids"`
}

type PermissionSpecModel struct {
	Role        types.String `tfsdk:"role"`
	Resource    types.String `tfsdk:"resource"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type NetworkSpecModel struct {
	Cidr           types.String `tfsdk:"cidr"`
	NetworkPortIds types.Set    `tfsdk:"network_port_ids"`
}

func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *policyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"service_type": schema.StringAttribute{
				Required: true,
			},
			"resource_ids": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"permission_specs": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							Required: true,
						},
						"resource": schema.StringAttribute{
							Required: true,
						},
						"permissions": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"network_spec": schema.SingleNestedAttribute{
				Optional: true,

				CustomType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"cidr": types.StringType,
						"network_port_ids": types.SetType{
							ElemType: types.StringType,
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"cidr": schema.StringAttribute{
						Required: true,
					},
					"network_port_ids": schema.SetAttribute{
						Required:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	rolesReq := make([]customer_metadata.MdsPermissionSpec, len(plan.PermissionSpecs))

	// Generate API request body from plan
	policyRequest := customer_metadata.MdsCreateUpdatePolicyRequest{
		Name:        plan.Name.ValueString(),
		ServiceType: plan.ServiceType.ValueString(),
	}
	if plan.ServiceType.ValueString() == policy_type.RABBITMQ {
		for i, roleId := range plan.PermissionSpecs {

			rolesReq[i] = customer_metadata.MdsPermissionSpec{
				Role:     roleId.Role.ValueString(),
				Resource: roleId.Resource.ValueString(),
			}
			roleId.Permissions.ElementsAs(ctx, &rolesReq[i].Permissions, true)
		}
		policyRequest.PermissionsSpec = rolesReq
	} else {
		networkSpec := &customer_metadata.MdsNetworkSpec{
			Cidr: plan.NetworkSpec.Cidr.ValueString(),
		}
		tflog.Debug(ctx, "Create Policy DTO", map[string]interface{}{"dto": plan.NetworkSpec.NetworkPortIds})

		plan.NetworkSpec.NetworkPortIds.ElementsAs(ctx, &networkSpec.NetworkPortIds, true)
		policyRequest.NetworkSpecs = append(policyRequest.NetworkSpecs, networkSpec)
	}

	tflog.Debug(ctx, "Create Policy DTO", map[string]interface{}{"dto": policyRequest})
	if _, err := r.client.CustomerMetadata.CreatePolicy(&policyRequest); err != nil {

		resp.Diagnostics.AddError(
			"Submitting request to create Policy",
			"Could not create Policy, unexpected error: "+err.Error(),
		)
		return
	}

	policies, err := r.client.CustomerMetadata.GetPolicies(&customer_metadata.MdsPoliciesQuery{
		Type:  plan.ServiceType.ValueString(),
		Names: []string{plan.Name.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError("Fetching Policy",
			"Could not fetch policy, unexpected error: "+err.Error(),
		)
		return
	}

	if len(*policies.Get()) <= 0 {
		resp.Diagnostics.AddError("Fetching Policy",
			"Unable to fetch the created policy",
		)
		return
	}
	createdPolicy := &(*policies.Get())[0]
	tflog.Debug(ctx, "Created Policy DTO", map[string]interface{}{"dto": createdPolicy})
	if saveFromPolicyResponse(&ctx, &resp.Diagnostics, &plan, createdPolicy) != 0 {
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

func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRequest := customer_metadata.MdsCreateUpdatePolicyRequest{
		Name:        plan.Name.ValueString(),
		ServiceType: plan.ServiceType.ValueString(),
	}
	rolesReq := make([]customer_metadata.MdsPermissionSpec, len(plan.PermissionSpecs))

	if plan.ServiceType.ValueString() == policy_type.RABBITMQ {
		for i, roleId := range plan.PermissionSpecs {
			rolesReq[i] = customer_metadata.MdsPermissionSpec{
				Role:     roleId.Role.ValueString(),
				Resource: roleId.Resource.ValueString(),
			}
			roleId.Permissions.ElementsAs(ctx, &rolesReq[i].Permissions, true)
		}
		updateRequest.PermissionsSpec = rolesReq
	} else {
		networkSpec := &customer_metadata.MdsNetworkSpec{
			Cidr: plan.NetworkSpec.Cidr.ValueString(),
		}
		plan.NetworkSpec.NetworkPortIds.ElementsAs(ctx, &networkSpec.NetworkPortIds, true)
		updateRequest.NetworkSpecs = append(updateRequest.NetworkSpecs, networkSpec)
	}
	tflog.Debug(ctx, "update policy request dto", map[string]interface{}{"dto": updateRequest})

	// Update existing policy
	if err := r.client.CustomerMetadata.UpdateMdsPolicy(plan.ID.ValueString(), &updateRequest); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS Policy",
			"Could not update Policy, unexpected error: "+err.Error(),
		)
		return
	}

	policy, err := r.client.CustomerMetadata.GetMDSPolicy(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching Policy",
			"Could not fetch policy while updating, unexpected error: "+err.Error(),
		)
		return
	}

	//Update resource state with updated items and timestamp
	if saveFromPolicyResponse(&ctx, &resp.Diagnostics, &plan, policy) != 0 {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *policyResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state policyResourceModel
	diags := request.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Submit request to delete MDS Policy
	err := r.client.CustomerMetadata.DeleteMdsPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting MDS Policy",
			"Could not delete MDS Policy by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed policy value from MDS
	policy, err := r.client.CustomerMetadata.GetMDSPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS Policy",
			"Could not read MDS policy ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	if saveFromPolicyResponse(&ctx, &resp.Diagnostics, &state, policy) != 0 {
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

func saveFromPolicyResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *policyResourceModel, policy *model.MdsPolicy) int8 {
	tflog.Info(*ctx, "Saving response to resourceModel state/plan", map[string]interface{}{"policy": *policy})
	if policy.ServiceType == service_type.RABBITMQ {
		roles, diags := convertFromPermissionSpecDto(ctx, policy.PermissionsSpec)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}
		state.PermissionSpecs = roles
	} else {
		tfNetworkSpecModels := make([]*NetworkSpecModel, len(policy.NetworkSpec))
		for i, networkSpec := range policy.NetworkSpec {
			tfNetworkSpecModels[i] = &NetworkSpecModel{
				Cidr: types.StringValue(networkSpec.CIDR),
			}
			networkPortIds, _ := types.SetValueFrom(*ctx, types.StringType, networkSpec.NetworkPortIds)
			tfNetworkSpecModels[i].NetworkPortIds = networkPortIds
		}
		state.NetworkSpec = tfNetworkSpecModels[0]
	}

	state.ID = types.StringValue(policy.ID)
	state.Name = types.StringValue(policy.Name)
	resourceIds, diags := types.SetValueFrom(*ctx, types.StringType, policy.ResourceIds)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.ResourceIds = resourceIds

	return 0
}

func convertFromPermissionSpecDto(ctx *context.Context, roles []*model.MdsPermissionsSpec) ([]*PermissionSpecModel, diag.Diagnostics) {
	tfPermissionSpecModels := make([]*PermissionSpecModel, len(roles))
	for i, permissionSpec := range roles {
		tfPermissionSpecModels[i] = &PermissionSpecModel{
			Role:     types.StringValue(permissionSpec.Role),
			Resource: types.StringValue(permissionSpec.Resource),
		}
		permissions, _ := types.SetValueFrom(*ctx, types.StringType, permissionSpec.Permissions)
		tfPermissionSpecModels[i].Permissions = permissions
	}

	return tfPermissionSpecModels, nil
}

func (r *policyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan policyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (plan.NetworkSpec == nil && len(plan.PermissionSpecs) > 0 && plan.ServiceType.ValueString() == policy_type.RABBITMQ) ||
		plan.NetworkSpec != nil && len(plan.PermissionSpecs) <= 0 && plan.ServiceType.ValueString() == policy_type.NETWORK {
		return
	}

	if plan.ServiceType.ValueString() == policy_type.RABBITMQ {
		resp.Diagnostics.AddError("ERROR:", "Cannot pass NetworkSpec as an input while creating RabbitMQ Policy.")
	} else {
		resp.Diagnostics.AddError("ERROR:", "Cannot pass PermissionSpecs as an input while creating a Network Policy.")
	}

}
