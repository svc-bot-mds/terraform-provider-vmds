package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
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
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
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
	state.Id = types.StringValue("placeholder")
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
