package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

var (
	_ datasource.DataSource              = &cloudAccountsDatasource{}
	_ datasource.DataSourceWithConfigure = &cloudAccountsDatasource{}
)

// instanceTypesDataSourceModel maps the data source schema data.
type cloudAccountsDatasourceModel struct {
	Id            types.String        `tfsdk:"id"`
	CloudAccounts []cloudAccountModel `tfsdk:"cloud_accounts"`
}

type cloudAccountModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AccountType types.String `tfsdk:"account_type"`
	Email       types.String `tfsdk:"user_email"`
	OrgId       types.String `tfsdk:"org_id"`
}

// NewCloudAccountsDatasource is a helper function to simplify the provider implementation.
func NewCloudAccountsDatasource() datasource.DataSource {
	return &cloudAccountsDatasource{}
}

// cloudAccountsDatasource is the data source implementation.
type cloudAccountsDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *cloudAccountsDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_accounts"
}

// Schema defines the schema for the data source.
func (d *cloudAccountsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch all cloud accounts on MDS for BYOC.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"cloud_accounts": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the cloud account.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the cloud account.",
							Computed:    true,
						},
						"account_type": schema.StringAttribute{
							Description: "Account type of the cloud account",
							Computed:    true,
						},
						"user_email": schema.StringAttribute{
							Description: "User email of the cloud account",
							Computed:    true,
						},
						"org_id": schema.StringAttribute{
							Description: "OrgId of the cloud account",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *cloudAccountsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudAccountsDatasourceModel
	var cloudAccountList []cloudAccountModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &infra_connector.MdsCloudAccountsQuery{}

	cloudAccounts, err := d.client.InfraConnector.GetCloudAccounts(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Cloud Accounts",
			err.Error(),
		)
		return
	}

	if cloudAccounts.Page.TotalPages > 1 {
		for i := 1; i <= cloudAccounts.Page.TotalPages; i++ {
			query.PageQuery.Index = i - 1
			totalCloudAccounts, err := d.client.InfraConnector.GetCloudAccounts(query)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Read MDS Cloud Accounts",
					err.Error(),
				)
				return
			}

			for _, cloudAccountDto := range *totalCloudAccounts.Get() {
				cloudAccount := cloudAccountModel{
					ID:          types.StringValue(cloudAccountDto.Id),
					Name:        types.StringValue(cloudAccountDto.Name),
					AccountType: types.StringValue(cloudAccountDto.AccountType),
					Email:       types.StringValue(cloudAccountDto.Email),
					OrgId:       types.StringValue(cloudAccountDto.OrgId),
				}
				cloudAccountList = append(cloudAccountList, cloudAccount)
			}
		}

		tflog.Debug(ctx, "cloud accounts dto", map[string]interface{}{"dto": cloudAccountList})
		state.CloudAccounts = append(state.CloudAccounts, cloudAccountList...)
	} else {
		for _, cloudAccountDto := range *cloudAccounts.Get() {
			tflog.Info(ctx, "Converting cloud account dto")
			cloudAccount := cloudAccountModel{
				ID:          types.StringValue(cloudAccountDto.Id),
				Name:        types.StringValue(cloudAccountDto.Name),
				AccountType: types.StringValue(cloudAccountDto.AccountType),
				Email:       types.StringValue(cloudAccountDto.Email),
				OrgId:       types.StringValue(cloudAccountDto.OrgId),
			}
			tflog.Debug(ctx, "converted cloud Account dto", map[string]interface{}{"dto": cloudAccount})
			state.CloudAccounts = append(state.CloudAccounts, cloudAccount)
		}
	}
	state.Id = types.StringValue(common.DataSource + common.CloudAccountsId)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *cloudAccountsDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
