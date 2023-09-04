package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

var (
	_ datasource.DataSource              = &serviceAccountsDatasource{}
	_ datasource.DataSourceWithConfigure = &serviceAccountsDatasource{}
)

// instanceTypesDataSourceModel maps the data source schema data.
type serviceAccountsDatasourceModel struct {
	Id              types.String          `tfsdk:"id"`
	ServiceAccounts []serviceAccountModel `tfsdk:"service_accounts"`
}

type serviceAccountModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
}

// NewServiceAccountsDataSource is a helper function to simplify the provider implementation.

func NewServiceAccountsDataSource() datasource.DataSource {
	return &serviceAccountsDatasource{}
}

// serviceAccountsDatasource is the data source implementation.
type serviceAccountsDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *serviceAccountsDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_accounts"
}

// Schema defines the schema for the data source.
func (d *serviceAccountsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch all service accounts on MDS for an Org (determined by the token used for provider).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"service_accounts": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the service account.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the service account.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the service account.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *serviceAccountsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serviceAccountsDatasourceModel
	var serviceAccountList []serviceAccountModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &customer_metadata.MdsServiceAccountsQuery{}

	serviceAccounts, err := d.client.CustomerMetadata.GetMdsServiceAccounts(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Service Accounts",
			err.Error(),
		)
		return
	}

	if serviceAccounts.Page.TotalPages > 1 {
		for i := 1; i <= serviceAccounts.Page.TotalPages; i++ {
			query.PageQuery.Index = i - 1
			totalServiceAccounts, err := d.client.CustomerMetadata.GetMdsServiceAccounts(query)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Read MDS Service Accounts",
					err.Error(),
				)
				return
			}

			for _, serviceAccountDto := range *totalServiceAccounts.Get() {
				serviceAccount := serviceAccountModel{
					ID:     types.StringValue(serviceAccountDto.Id),
					Name:   types.StringValue(serviceAccountDto.Name),
					Status: types.StringValue(serviceAccountDto.Status),
				}
				serviceAccountList = append(serviceAccountList, serviceAccount)
			}
		}

		tflog.Debug(ctx, "service accounts dto", map[string]interface{}{"dto": serviceAccountList})
		state.ServiceAccounts = append(state.ServiceAccounts, serviceAccountList...)
	} else {
		for _, serviceAccountDto := range *serviceAccounts.Get() {
			tflog.Info(ctx, "Converting svc account dto")
			serviceAccount := serviceAccountModel{
				ID:     types.StringValue(serviceAccountDto.Id),
				Name:   types.StringValue(serviceAccountDto.Name),
				Status: types.StringValue(serviceAccountDto.Status),
			}
			tflog.Debug(ctx, "converted service Account dto", map[string]interface{}{"dto": serviceAccount})
			state.ServiceAccounts = append(state.ServiceAccounts, serviceAccount)
		}
	}
	state.Id = types.StringValue(common.DataSource + common.ServiceAccountsId)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceAccountsDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
