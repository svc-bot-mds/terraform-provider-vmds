package mds

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/time_unit"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
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
	ID         types.String          `tfsdk:"id"`
	Name       types.String          `tfsdk:"name"`
	Status     types.String          `tfsdk:"status"`
	PolicyIds  types.Set             `tfsdk:"policy_ids"`
	Tags       types.Set             `tfsdk:"tags"`
	Timeouts   timeouts.Value        `tfsdk:"timeouts"`
	Credential types.Object          `tfsdk:"credential"`
	OauthApp   basetypes.ObjectValue `tfsdk:"oauth_app"`
}

type ServiceAccountCredential struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	OrgId        types.String `tfsdk:"org_id"`
	GrantType    types.String `tfsdk:"grant_type"`
}
type ServiceAccountOauthApp struct {
	OauthAppId  types.String `tfsdk:"app_id"`
	AppType     types.String `tfsdk:"app_type"`
	Created     types.String `tfsdk:"created"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Description types.String `tfsdk:"description"`
	Modified    types.String `tfsdk:"modified"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	TTLSpec     TTLSpecModel `tfsdk:"ttl_spec"`
}

type TTLSpecModel struct {
	TimeUnit types.String `tfsdk:"time_unit"`
	TTL      types.Int64  `tfsdk:"ttl"`
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
		MarkdownDescription: "Represents a service account created on MDS, can be used to create/update/delete/import a service account.\n" +
			"Note: 1. Only service accounts with valid oAuthapp can be imported.\n2. Please make sure you have selected the valid policy with active clusters while creating the service account",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Auto-generated ID after creating an user, and can be passed to import an existing user from MDS to terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Updating the name results in deletion of existing service account and new service account with updated name is created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description: "Active status of service account on MDS.",
				Computed:    true,
			},
			"policy_ids": schema.SetAttribute{
				Description: "IDs of service policies to be associated with service account.",
				Optional:    true,
				Computed:    false,
				ElementType: types.StringType,
			},
			"tags": schema.SetAttribute{
				Description: "Tags or labels to categorise service accounts for ease of finding.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"oauth_app": schema.SingleNestedAttribute{
				MarkdownDescription: "Provides OauthApp details.",
				Computed:            true,
				Optional:            true,
				CustomType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"app_id":      types.StringType,
						"app_type":    types.StringType,
						"created":     types.StringType,
						"created_by":  types.StringType,
						"description": types.StringType,
						"modified":    types.StringType,
						"modified_by": types.StringType,
						"ttl_spec": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"ttl":       types.Int64Type,
								"time_unit": types.StringType,
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"app_id": schema.StringAttribute{
						Description: "Id of the oAuthApp.",
						Computed:    true,
					},
					"app_type": schema.StringAttribute{
						Description: "Type of the oAuthApp.",
						Computed:    true,
					},
					"created": schema.StringAttribute{
						Description: "Time when the service account was created.",
						Computed:    true,
					},
					"created_by": schema.StringAttribute{
						Description: "Username of the user who has created the service account.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "Description of the OauthApp.",
						Computed:    true,
						Optional:    true,
					},
					"modified_by": schema.StringAttribute{
						Description: "Username of the user who has updated the service account.",
						Computed:    true,
					},
					"modified": schema.StringAttribute{
						Description: "Time when the service account was modified.",
						Computed:    true,
					},
					"ttl_spec": schema.SingleNestedAttribute{
						Description: "OauthApp Access token Duration details. Valid TTL value : less than 5 hours or 300 minutes.",
						Computed:    true,
						Optional:    true,
						CustomType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"ttl":       types.Int64Type,
								"time_unit": types.StringType,
							},
						},
						Attributes: map[string]schema.Attribute{
							"time_unit": schema.StringAttribute{
								MarkdownDescription: "Unit of time. Valid values : `HOURS` or `MINUTES`.",
								Computed:            true,
								Optional:            true,
							},
							"ttl": schema.Int64Attribute{
								Description: "time to live value.",
								Computed:    true,
								Optional:    true,
							},
						},
					},
				},
			},
			"credential": schema.SingleNestedAttribute{
				MarkdownDescription: "Holds the Client Secret details.",
				Computed:            true,
				Optional:            true,
				CustomType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"client_id":     types.StringType,
						"client_secret": types.StringType,
						"grant_type":    types.StringType,
						"org_id":        types.StringType,
					},
				},
				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: "Client Id generated for the service account.",
						Computed:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Client Secret generated for the service account.",
						Computed:    true,
					},
					"grant_type": schema.StringAttribute{
						Description: "Grant Type of the credentials.",
						Computed:    true,
					},
					"org_id": schema.StringAttribute{
						Description: "Org Id of the current user.",
						Computed:    true,
					},
				},
			},
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
		Usernames: []string{plan.Name.ValueString()},
	}
	plan.Tags.ElementsAs(ctx, &svcAccountRequest.Tags, true)
	plan.PolicyIds.ElementsAs(ctx, &svcAccountRequest.PolicyIds, true)

	svcAcctCredentials, err := r.client.CustomerMetadata.CreateMdsServiceAccount(&svcAccountRequest)
	if err != nil {
		apiErr := core.ApiError{}
		errors.As(err, &apiErr)
		if apiErr.ErrorCode == customer_metadata.DuplicateServiceAccount {
			resp.Diagnostics.AddError(
				"Submitting request to create service account",
				"There was some issue while creating the service account."+
					" Unexpected error: "+err.Error(),
			)
		} else {
			resp.Diagnostics.AddError(
				"Submitting request to create service account",
				"There was some issue while creating the service account and oauth app for the service account ."+
					" Please verify if the service account was created without an oauth app and delete it."+
					" Unexpected error: "+err.Error(),
			)
		}
		return
	}

	svcAccounts, err := r.client.CustomerMetadata.GetMdsServiceAccounts(&customer_metadata.MdsServiceAccountsQuery{
		Names: []string{plan.Name.ValueString()},
	})
	createdSvcAcct := &(*svcAccounts.Get())[0]

	svcAccountsOauthApps, oauthError := r.client.CustomerMetadata.GetMDSServiceAccountOauthApp(createdSvcAcct.Id)

	if oauthError != nil {
		resp.Diagnostics.AddError("Fetching oAuth Apps for the Service Account",
			"Could not fetch oAuth Apps for the Service Account, unexpected error: "+err.Error(),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Fetching service account",
			"Could not fetch service accounts, unexpected error: "+err.Error(),
		)
		return
	}
	if svcAccounts.Page.TotalElements == 0 {
		resp.Diagnostics.AddError("Fetching Service Account",
			fmt.Sprintf("Could not find any Service Account by name [%s], server error must have occurred while creating svc account.", plan.Name.ValueString()),
		)
		return
	}

	if len(*svcAccounts.Get()) <= 0 {
		resp.Diagnostics.AddError("Fetching Service Accounts",
			"Unable to fetch the created service account",
		)
		return
	}

	var serviceAccountOauthAppPlan = convertToOauthAppModel(ctx, &plan)
	if serviceAccountOauthAppPlan.TTLSpec.TTL.ValueInt64() != svcAccountsOauthApps.TTLSpec.TTL ||
		serviceAccountOauthAppPlan.Description.ValueString() != svcAccountsOauthApps.Description ||
		serviceAccountOauthAppPlan.TTLSpec.TimeUnit.ValueString() != svcAccountsOauthApps.TTLSpec.TimeUnit {
		updateRequest := customer_metadata.MDSOauthAppUpdateRequest{
			Description: serviceAccountOauthAppPlan.Description.ValueString(),
			TTL:         serviceAccountOauthAppPlan.TTLSpec.TTL.ValueInt64(),
			TimeUnit:    serviceAccountOauthAppPlan.TTLSpec.TimeUnit.ValueString(),
		}
		svcAccountsOauthApps, err = r.client.CustomerMetadata.UpdateMDSServiceAccountOauthApp(createdSvcAcct.Id, &updateRequest, svcAccountsOauthApps.AppId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Creating MDS service account - Oauth App details",
				"Could not update service account - oauth app details, unexpected error: "+err.Error(),
			)
			return
		}

	}
	// Map response body to schema and populate Computed attribute values
	if saveFromSvcAccountCreateResponse(&ctx, &resp.Diagnostics, &plan, svcAcctCredentials, createdSvcAcct, svcAccountsOauthApps) != 0 {
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
	var state serviceAccountResourceModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRequest := customer_metadata.MdsSvcAccountUpdateRequest{}
	state.Tags.ElementsAs(ctx, &updateRequest.Tags, true)
	state.PolicyIds.ElementsAs(ctx, &updateRequest.PolicyIds, true)

	// Update existing svc account
	if err := r.client.CustomerMetadata.UpdateMdsServiceAccount(state.ID.ValueString(), &updateRequest); err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS service account",
			"Could not update service account, unexpected error: "+err.Error(),
		)
		return
	}

	svcAccount, err := r.client.CustomerMetadata.GetMdsServiceAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Fetching svc account",
			"Could not fetch svc account while updating, unexpected error: "+err.Error(),
		)
		return
	}
	svcAccountsOauthAppResponse, oauthError := r.client.CustomerMetadata.GetMDSServiceAccountOauthApp(state.ID.ValueString())
	if oauthError != nil {
		resp.Diagnostics.AddError("Fetching oAuth Apps for the Service Account",
			"Could not fetch oAuth Apps for the Service Account, unexpected error: "+err.Error(),
		)
		return
	}
	var serviceAccountOauthApp = convertToOauthAppModel(ctx, &state)
	var oauthApp *model.MDSServieAccountOauthApp
	if serviceAccountOauthApp.Description.ValueString() != "" || serviceAccountOauthApp.TTLSpec.TTL.ValueInt64() != 0 ||
		serviceAccountOauthApp.TTLSpec.TimeUnit.ValueString() != "" {
		updateRequest := customer_metadata.MDSOauthAppUpdateRequest{
			Description: serviceAccountOauthApp.Description.ValueString(),
			TTL:         serviceAccountOauthApp.TTLSpec.TTL.ValueInt64(),
			TimeUnit:    serviceAccountOauthApp.TTLSpec.TimeUnit.ValueString(),
		}
		oauthApp, err = r.client.CustomerMetadata.UpdateMDSServiceAccountOauthApp(state.ID.ValueString(), &updateRequest, svcAccountsOauthAppResponse.AppId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Updating MDS service account - Oauth App details",
				"Could not update service account - oauth app details, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if saveFromSvcAccountCreateResponseFromUpdate(&ctx, &resp.Diagnostics, &state, svcAccount, oauthApp) != 0 {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "END__Update")
}

func convertToOauthAppModel(ctx context.Context, state *serviceAccountResourceModel) *ServiceAccountOauthApp {
	var serviceAccountOauthApp = &ServiceAccountOauthApp{}
	var options = basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	}
	state.OauthApp.As(ctx, serviceAccountOauthApp, options)
	return serviceAccountOauthApp
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

	// Get refreshed service account value from MDS
	svcAcct, err := r.client.CustomerMetadata.GetMdsServiceAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Updating MDS service account",
			"Could not update service account, unexpected error: "+err.Error(),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Reading MDS service account",
			"Could not read MDS service account ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	svcAccountsOauthApps, oauthError := r.client.CustomerMetadata.GetMDSServiceAccountOauthApp(state.ID.ValueString())

	if oauthError != nil {
		resp.Diagnostics.AddError("Fetching oAuth Apps for the Service Account",
			"Could not fetch oAuth Apps for the Service Account, unexpected error: "+err.Error(),
		)
		return
	}
	// Overwrite items with refreshed state

	if saveFromSvcAccountCreateResponse(&ctx, &resp.Diagnostics, &state, nil, svcAcct, svcAccountsOauthApps) != 0 {
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

func saveFromSvcAccountCreateResponse(ctx *context.Context, diagnostics *diag.Diagnostics, state *serviceAccountResourceModel,
	serviceAccountCreateResponse *model.MdsServiceAccountCreate, createdSvcAccount *model.MdsServiceAccount, svcAccountsOauthApps *model.MDSServieAccountOauthApp) int8 {
	state.Name = types.StringValue(createdSvcAccount.Name)
	state.Status = types.StringValue(createdSvcAccount.Status)
	tags, diags := types.SetValueFrom(*ctx, types.StringType, createdSvcAccount.Tags)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.Tags = tags
	state.ID = types.StringValue(createdSvcAccount.Id)
	//Client Secret
	if serviceAccountCreateResponse != nil && serviceAccountCreateResponse.OAuthCredentials != nil && serviceAccountCreateResponse.OAuthCredentials[0].Credential != nil {

		credentialModel := ServiceAccountCredential{
			ClientSecret: types.StringValue(serviceAccountCreateResponse.OAuthCredentials[0].Credential.ClientSecret),
			ClientId:     types.StringValue(serviceAccountCreateResponse.OAuthCredentials[0].Credential.ClientId),
			OrgId:        types.StringValue(serviceAccountCreateResponse.OAuthCredentials[0].Credential.OrgId),
			GrantType:    types.StringValue(serviceAccountCreateResponse.OAuthCredentials[0].Credential.GrantType),
		}
		credentialsObject, diags := types.ObjectValueFrom(*ctx, state.Credential.AttributeTypes(*ctx), credentialModel)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}
		state.Credential = credentialsObject

	}

	//oauth App
	if svcAccountsOauthApps != nil {
		oauthAppModel := ServiceAccountOauthApp{
			OauthAppId:  types.StringValue(svcAccountsOauthApps.AppId),
			AppType:     types.StringValue(svcAccountsOauthApps.AppType),
			Created:     types.StringValue(svcAccountsOauthApps.Created),
			CreatedBy:   types.StringValue(svcAccountsOauthApps.CreatedBy),
			Description: types.StringValue(svcAccountsOauthApps.Description),
			Modified:    types.StringValue(svcAccountsOauthApps.Modified),
			ModifiedBy:  types.StringValue(svcAccountsOauthApps.ModifiedBy),
		}
		ttlSpecModel := TTLSpecModel{
			TimeUnit: types.StringValue(svcAccountsOauthApps.TTLSpec.TimeUnit),
			TTL:      types.Int64Value(svcAccountsOauthApps.TTLSpec.TTL),
		}

		oauthAppModel.TTLSpec = ttlSpecModel
		oauthObject, diags := types.ObjectValueFrom(*ctx, state.OauthApp.AttributeTypes(*ctx), oauthAppModel)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}

		state.OauthApp = oauthObject
	}
	return 0
}

func saveFromSvcAccountCreateResponseFromUpdate(ctx *context.Context, diagnostics *diag.Diagnostics, state *serviceAccountResourceModel, createdSvcAccount *model.MdsServiceAccount, svcAccountsOauthApps *model.MDSServieAccountOauthApp) int8 {

	tflog.Info(*ctx, "Saving create response to resourceModel state/plan", map[string]interface{}{"service_accounts": *createdSvcAccount})

	state.Name = types.StringValue(createdSvcAccount.Name)
	state.Status = types.StringValue(createdSvcAccount.Status)
	tags, diags := types.SetValueFrom(*ctx, types.StringType, createdSvcAccount.Tags)
	if diagnostics.Append(diags...); diagnostics.HasError() {
		return 1
	}
	state.Tags = tags
	state.ID = types.StringValue(createdSvcAccount.Id)

	//Client Secret
	if state.Credential.IsUnknown() {
		credentialModel := ServiceAccountCredential{
			ClientSecret: types.StringValue(""),
			ClientId:     types.StringValue(""),
			OrgId:        types.StringValue(""),
			GrantType:    types.StringValue(""),
		}
		credentialsObject, diags := types.ObjectValueFrom(*ctx, state.Credential.AttributeTypes(*ctx), credentialModel)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}
		state.Credential = credentialsObject

	}

	//oauth App
	if svcAccountsOauthApps != nil {
		oauthAppModel := ServiceAccountOauthApp{
			OauthAppId:  types.StringValue(svcAccountsOauthApps.AppId),
			AppType:     types.StringValue(svcAccountsOauthApps.AppType),
			Created:     types.StringValue(svcAccountsOauthApps.Created),
			CreatedBy:   types.StringValue(svcAccountsOauthApps.CreatedBy),
			Description: types.StringValue(svcAccountsOauthApps.Description),
			Modified:    types.StringValue(svcAccountsOauthApps.Modified),
			ModifiedBy:  types.StringValue(svcAccountsOauthApps.ModifiedBy),
		}
		ttlSpecModel := TTLSpecModel{
			TimeUnit: types.StringValue(svcAccountsOauthApps.TTLSpec.TimeUnit),
			TTL:      types.Int64Value(svcAccountsOauthApps.TTLSpec.TTL),
		}

		oauthAppModel.TTLSpec = ttlSpecModel
		oauthObject, diags := types.ObjectValueFrom(*ctx, state.OauthApp.AttributeTypes(*ctx), oauthAppModel)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}

		state.OauthApp = oauthObject
	}
	return 0
}

func (r *serviceAccountResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan serviceAccountResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.OauthApp.IsNull() {
		oauthAppModel := convertToOauthAppModel(ctx, &plan)
		if (oauthAppModel.TTLSpec.TTL.ValueInt64() > 300 && oauthAppModel.TTLSpec.TimeUnit.ValueString() == time_unit.MINUTES) ||
			(oauthAppModel.TTLSpec.TTL.ValueInt64() > 5 && oauthAppModel.TTLSpec.TimeUnit.ValueString() == time_unit.HOURS) ||
			(oauthAppModel.TTLSpec.TimeUnit.ValueString() != time_unit.MINUTES && oauthAppModel.TTLSpec.TimeUnit.ValueString() != time_unit.HOURS) {
			resp.Diagnostics.AddError("Validation Failed", "Please enter a valid TTL value less than 5 hours or 300 minutes.")
		}
	}
}
