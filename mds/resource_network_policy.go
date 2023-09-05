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
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"regexp"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkPolicyResource{}
	_ resource.ResourceWithConfigure   = &networkPolicyResource{}
	_ resource.ResourceWithImportState = &networkPolicyResource{}
)

func NewNetworkPolicyResource() resource.Resource {
	return &networkPolicyResource{}
}

type networkPolicyResource struct {
	client *mds.Client
}

type networkPolicyResourceModel struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	ServiceType types.String      `tfsdk:"service_type"`
	NetworkSpec *NetworkSpecModel `tfsdk:"network_spec"`
	ResourceIds types.Set         `tfsdk:"resource_ids"`
}

type NetworkSpecModel struct {
	Cidr           types.String `tfsdk:"cidr"`
	NetworkPortIds types.Set    `tfsdk:"network_port_ids"`
}

func (r *networkPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_policy"
}

func (r *networkPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *networkPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		Description: "Represents a policy on MDS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID of the policy after creation, and can be used to import it from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_type": schema.StringAttribute{
				MarkdownDescription: "Type of policy to manage. Supported values is:  `NETWORK`.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy",
				Required:    true,
			},
			"resource_ids": schema.SetAttribute{
				Description: "IDs of service resources/instances being managed by the policy.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"network_spec": schema.SingleNestedAttribute{
				MarkdownDescription: "Network config to allow access to service resource.",
				Required:            true,
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
						MarkdownDescription: "CIDR value to allow access from. Ex: `10.45.66.80/30`",
						Required:            true,
					},
					"network_port_ids": schema.SetAttribute{
						MarkdownDescription: "IDs of network ports to open up for access. Please make use of datasource `vmds_network_ports` to get IDs of ports available for services.",
						Required:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *networkPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan networkPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	policyRequest := customer_metadata.MdsCreateUpdatePolicyRequest{
		Name:        plan.Name.ValueString(),
		ServiceType: policy_type.NETWORK,
	}

	networkSpec := &customer_metadata.MdsNetworkSpec{
		Cidr: plan.NetworkSpec.Cidr.ValueString(),
	}
	tflog.Debug(ctx, "Create Network Policy DTO", map[string]interface{}{"dto": plan.NetworkSpec.NetworkPortIds})

	plan.NetworkSpec.NetworkPortIds.ElementsAs(ctx, &networkSpec.NetworkPortIds, true)
	policyRequest.NetworkSpecs = append(policyRequest.NetworkSpecs, networkSpec)

	tflog.Debug(ctx, "Create Network Policy DTO", map[string]interface{}{"dto": policyRequest})
	if _, err := r.client.CustomerMetadata.CreatePolicy(&policyRequest); err != nil {

		resp.Diagnostics.AddError(
			"Submitting request to create Network Policy",
			"Could not create Network Policy, unexpected error: "+err.Error(),
		)
		return
	}

	policies, err := r.client.CustomerMetadata.GetPolicies(&customer_metadata.MdsPoliciesQuery{
		Type:  policy_type.NETWORK,
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
	if saveFromNetworkPolicyResponse(&ctx, &resp.Diagnostics, &plan, createdPolicy) != 0 {
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

func (r *networkPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var plan networkPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRequest := customer_metadata.MdsCreateUpdatePolicyRequest{
		Name:        plan.Name.ValueString(),
		ServiceType: policy_type.NETWORK,
	}

	networkSpec := &customer_metadata.MdsNetworkSpec{
		Cidr: plan.NetworkSpec.Cidr.ValueString(),
	}
	plan.NetworkSpec.NetworkPortIds.ElementsAs(ctx, &networkSpec.NetworkPortIds, true)
	updateRequest.NetworkSpecs = append(updateRequest.NetworkSpecs, networkSpec)

	tflog.Debug(ctx, "update policy request dto", map[string]interface{}{"dto": updateRequest})

	// Update existing policy
	if err := r.client.CustomerMetadata.UpdateMdsPolicy(plan.ID.ValueString(), &updateRequest); err != nil {
		resp.Diagnostics.AddError(
			"Updating  Network Policy",
			"Could not update Network Policy, unexpected error: "+err.Error(),
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
	if saveFromNetworkPolicyResponse(&ctx, &resp.Diagnostics, &plan, policy) != 0 {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *networkPolicyResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state networkPolicyResourceModel
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

func (r *networkPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *networkPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state networkPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed policy value from MDS
	policy, err := r.client.CustomerMetadata.GetMDSPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading Network Policy",
			"Could not read Network policy ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	if saveFromNetworkPolicyResponse(&ctx, &resp.Diagnostics, &state, policy) != 0 {
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

func saveFromNetworkPolicyResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *networkPolicyResourceModel, policy *model.MdsPolicy) int8 {
	tflog.Info(*ctx, "Saving response to resourceModel state/plan", map[string]interface{}{"policy": *policy})

	tfNetworkSpecModels := make([]*NetworkSpecModel, len(policy.NetworkSpec))
	for i, networkSpec := range policy.NetworkSpec {
		tfNetworkSpecModels[i] = &NetworkSpecModel{
			Cidr: types.StringValue(networkSpec.CIDR),
		}
		networkPortIds, _ := types.SetValueFrom(*ctx, types.StringType, networkSpec.NetworkPortIds)
		tfNetworkSpecModels[i].NetworkPortIds = networkPortIds
	}
	state.NetworkSpec = tfNetworkSpecModels[0]

	state.ID = types.StringValue(policy.ID)
	state.Name = types.StringValue(policy.Name)
	resourceIds, diags := types.SetValueFrom(*ctx, types.StringType, policy.ResourceIds)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.ResourceIds = resourceIds

	return 0
}

func (r *networkPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan networkPolicyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.NetworkSpec != nil {
		const pattern = `^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))$`
		regex := regexp.MustCompile(pattern)

		// Check if the value matches the pattern
		if !regex.MatchString(plan.NetworkSpec.Cidr.ValueString()) {
			resp.Diagnostics.AddError("Validation Failed", "CIDR form is invalid.( Ex. 10.22.55.0/24 )")

		}
	} else {
		return
	}

}
