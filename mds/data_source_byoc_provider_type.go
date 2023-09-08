package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &providerTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &providerTypesDataSource{}
)

type providerTypesDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	ProviderTypes []string     `tfsdk:"provider_types"`
}

func NewProviderTypesDataSource() datasource.DataSource {
	return &providerTypesDataSource{}
}

type providerTypesDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *providerTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_types"
}

func (d *providerTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch types of providers supported by BYOC.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource.",
			},
			"provider_types": schema.SetAttribute{
				Description: "List of provider types.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *providerTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *providerTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "getProviderTypes")
	var state providerTypesDataSourceModel

	//Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	tflog.Info(ctx, "getProviderTypes")
	typesList, err := d.client.InfraConnector.GetProviderTypes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Provider Types:",
			err.Error(),
		)
		return
	}
	state.Id = types.StringValue(common.DataSource + common.ProviderTypesId)
	state.ProviderTypes = append(state.ProviderTypes, typesList...)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
