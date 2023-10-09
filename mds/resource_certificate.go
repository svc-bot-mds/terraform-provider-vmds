package mds

import (
	"context"
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
	"net/url"
	"time"
)

const (
	defaultCertificateCreateTimeout = 2 * time.Minute
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithConfigure   = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

type certificateResource struct {
	client *mds.Client
}

type CertificateResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	DomainName     types.String   `tfsdk:"domain_name"`
	ProviderType   types.String   `tfsdk:"provider_type"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	ExpirationTime types.String   `tfsdk:"expiration_time"`
	CreatedBy      types.String   `tfsdk:"created_by"`
	Certificate    types.String   `tfsdk:"certificate"`
	CertificateCA  types.String   `tfsdk:"certificate_ca"`
	CertificateKey types.String   `tfsdk:"certificate_key"`
}

func (r *certificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *certificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *certificateResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "INIT__Schema")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents a certificate created on MDS, can be used to create/update/delete/import a certificate.\n" +
			fmt.Sprintf("3. Default timeout for creation is `%v`.", defaultCertificateCreateTimeout),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID after creating a certificate, and can be passed to import an existing user from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name is readonly field while updating the certificate.",
				Required:    true,
			},
			"provider_type": schema.StringAttribute{
				Description: "Provider Type of certificate on MDS. It is a readonly field while updating the certificate.",
				Required:    true,
			},
			"domain_name": schema.StringAttribute{
				Description: "Domain Name of the certificate on MDS. It is a readonly field while updating the certificate.",
				Required:    true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"expiration_time": schema.StringAttribute{
				MarkdownDescription: "Holds the ExpirationTime of the certificate.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Email of the user who has created the certificate",
				Computed:            true,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "Certificate details",
				Required:            true,
			},
			"certificate_ca": schema.StringAttribute{
				MarkdownDescription: "Certificate CA details",
				Required:            true,
			},
			"certificate_key": schema.StringAttribute{
				MarkdownDescription: "Certificate Key details",
				Required:            true,
			},
		},
	}

	tflog.Info(ctx, "END__Schema")
}

// Create a new resource
func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "INIT__Create")
	// Retrieve values from plan
	var plan CertificateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCertificateCreateTimeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Generate API request body from plan
	certificateRequest := &infra_connector.CertificateCreateRequest{
		Name:           plan.Name.ValueString(),
		DomainName:     plan.DomainName.ValueString(),
		Provider:       plan.ProviderType.ValueString(),
		Certificate:    plan.Certificate.ValueString(),
		CertificateCA:  plan.CertificateCA.ValueString(),
		CertificateKey: plan.CertificateKey.ValueString(),
	}

	tflog.Info(ctx, "req param", map[string]interface{}{"reeed": certificateRequest})
	certificate, err := r.client.InfraConnector.CreateCertificate(certificateRequest)
	if err != nil {
		apiErr := core.ApiError{}
		errors.As(err, &apiErr)
		resp.Diagnostics.AddError(
			"Submitting request to create certificate",
			"There was some issue while creating the certificate."+
				" Unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	if saveFromCertificateCreateResponse(&plan, certificate) != 0 {
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

func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "INIT__Update")

	// Retrieve values from plan
	var state CertificateResourceModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificateUpdateReq := infra_connector.CertificateUpdateRequest{
		Certificate:    url.PathEscape(state.Certificate.ValueString()),
		CertificateCA:  url.PathEscape(state.CertificateCA.ValueString()),
		CertificateKey: url.PathEscape(state.CertificateKey.ValueString()),
	}

	// Update existing svc account
	if _, err := r.client.InfraConnector.UpdateCertificate(state.ID.ValueString(), &certificateUpdateReq); err != nil {
		resp.Diagnostics.AddError(
			"Updating the Certificate",
			"Could not update certificate, unexpected error: "+err.Error(),
		)
		return
	}

	certificate, err := r.client.InfraConnector.GetCertificate(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching certificate",
			"Could not fetch certificate while updating, unexpected error: "+err.Error(),
		)
		return
	}

	if saveFromCertificateCreateResponse(&state, &certificate) != 0 {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func (r *certificateResource) Delete(ctx context.Context, request resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "INIT__Delete")
	// Get current state
	var state CertificateResourceModel
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

	// Submit request to delete MDS certificate
	err := r.client.InfraConnector.DeleteCertificate(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Deleting certificate",
			"Could not delete certificate by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "END__Delete")
}

func (r *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "INIT__Read")
	// Get current state
	var state CertificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed certificate value from MDS
	certificate, err := r.client.InfraConnector.GetCertificate(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS certificate",
			"Could not read MDS certificate "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if saveFromCertificateCreateResponse(&state, &certificate) != 0 {
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

func saveFromCertificateCreateResponse(state *CertificateResourceModel,
	certificateResponse *model.MdsCertificate) int8 {
	state.Name = types.StringValue(certificateResponse.Name)
	state.DomainName = types.StringValue(certificateResponse.DomainName)
	state.ID = types.StringValue(certificateResponse.Id)
	state.ProviderType = types.StringValue(certificateResponse.Provider)
	state.ExpirationTime = types.StringValue(certificateResponse.ExpiryTime)
	state.CreatedBy = types.StringValue(certificateResponse.CreatedBy)
	return 0
}
