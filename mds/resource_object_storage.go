package mds

import (
	"context"
	"errors"
	"fmt"
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
	"net/url"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &objectStorageResource{}
	_ resource.ResourceWithConfigure   = &objectStorageResource{}
	_ resource.ResourceWithImportState = &objectStorageResource{}
)

func NewObjectStorageResource() resource.Resource {
	return &objectStorageResource{}
}

type objectStorageResource struct {
	client *mds.Client
}

type ObjectStorageResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	BucketName      types.String `tfsdk:"bucket_name"`
	Endpoint        types.String `tfsdk:"endpoint"`
	Region          types.String `tfsdk:"region"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
}

func (r *objectStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage"
}

func (r *objectStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *objectStorageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents a object storage created on MDS, can be used to create/update/delete/import a object storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID after creating a object storage, and can be passed to import an existing user from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the object storage. Updating it will result in creating new object store.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_name": schema.StringAttribute{
				Description: "Name of the initial bucket to create. Modifying this field is not allowed.",
				Required:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint of the object storage to use. Modifying this field is not allowed.",
				Required:    true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Region where object storage is created. Modifying this field is not allowed.",
				Required:            true,
			},
			"access_key_id": schema.StringAttribute{
				MarkdownDescription: "Access Key Id for the authentication of object storage.",
				Required:            true,
			},
			"secret_access_key": schema.StringAttribute{
				MarkdownDescription: "Secret Access Key for the authentication of object storage.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *objectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan ObjectStorageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	request := &infra_connector.ObjectStorageCreateRequest{
		Name:            plan.Name.ValueString(),
		BucketName:      plan.BucketName.ValueString(),
		Endpoint:        plan.Endpoint.ValueString(),
		Region:          plan.Region.ValueString(),
		AccessKeyId:     plan.AccessKeyId.ValueString(),
		SecretAccessKey: plan.SecretAccessKey.ValueString(),
	}

	tflog.Info(ctx, "req param", map[string]interface{}{"create-request": request})
	response, err := r.client.InfraConnector.CreateObjectStorage(request)
	if err != nil {
		apiErr := core.ApiError{}
		errors.As(err, &apiErr)
		resp.Diagnostics.AddError(
			"Submitting request to create object store",
			"There was some issue while creating the object store."+
				" Unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	if r.saveFromResponse(&plan, response) != 0 {
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

func (r *objectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var state, plan ObjectStorageResourceModel
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
		return
	}

	request := infra_connector.ObjectStorageUpdateRequest{
		AccessKeyId:     url.PathEscape(plan.AccessKeyId.ValueString()),
		SecretAccessKey: url.PathEscape(plan.SecretAccessKey.ValueString()),
	}

	// Update existing svc account
	if _, err := r.client.InfraConnector.UpdateObjectStore(plan.ID.ValueString(), &request); err != nil {
		resp.Diagnostics.AddError(
			"Updating the object store",
			"Could not update object store, unexpected error: "+err.Error(),
		)
		return
	}

	response, err := r.client.InfraConnector.GetObjectStorage(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching object store",
			"Could not fetch object store while updating, unexpected error: "+err.Error(),
		)
		return
	}

	if r.saveFromResponse(&state, &response) != 0 {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *objectStorageResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state ObjectStorageResourceModel
	diags := request.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Submit request to delete MDS certificate
	err := r.client.InfraConnector.DeleteObjectStorage(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting object store",
			"Could not delete object store by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *objectStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *objectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state ObjectStorageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed certificate value from MDS
	response, err := r.client.InfraConnector.GetObjectStorage(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS object store",
			"Could not read MDS object store "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if r.saveFromResponse(&state, &response) != 0 {
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

func (r *objectStorageResource) saveFromResponse(state *ObjectStorageResourceModel, response *model.MdsObjectStorage) int8 {
	state.ID = types.StringValue(response.Id)
	return 0
}

func (r *objectStorageResource) validateUpdateRequest(state *ObjectStorageResourceModel, plan *ObjectStorageResourceModel) string {
	if state.Region != plan.Region {
		return `Updating "region" is not allowed`
	}
	if state.Endpoint != plan.Endpoint {
		return `Updating "endpoint" is not allowed`
	}
	if state.BucketName != plan.BucketName {
		return `Updating "bucket_name" is not allowed`
	}
	return "OK"
}
