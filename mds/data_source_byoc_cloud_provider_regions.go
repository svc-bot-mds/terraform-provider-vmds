package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &cloudProviderRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &cloudProviderRegionsDataSource{}
)

type cloudProviderRegionsDataSourceModel struct {
	Id                   types.String               `tfsdk:"id"`
	CloudProviderRegions []cloudProviderRegionModel `tfsdk:"cloud_provider_regions"`
}

type cloudProviderRegionModel struct {
	Name      types.String `tfsdk:"name"`
	Regions   types.List   `tfsdk:"regions"`
	ShortName types.String `tfsdk:"short_name"`
	Id        types.String `tfsdk:"id"`
}

func NewCloudProviderRegionsDataSource() datasource.DataSource {
	return &cloudProviderRegionsDataSource{}
}

type cloudProviderRegionsDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *cloudProviderRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_provider_regions"
}

func (d *cloudProviderRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch the list of cloud providers with the regions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"cloud_provider_regions": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the cloud provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the cloud provider.",
							Computed:    true,
						},
						"short_name": schema.StringAttribute{
							Description: "Short Name of the cloud provider",
							Computed:    true,
						},
						"regions": schema.ListAttribute{
							Description: "List of Regions.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *cloudProviderRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *cloudProviderRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudProviderRegionsDataSourceModel

	//Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	regions, err := d.client.InfraConnector.GetCloudProviderRegions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Regions:",
			err.Error(),
		)
		return
	}
	state.Id = types.StringValue(common.DataSource + common.RegionsId)

	for _, cloudProvider := range regions {
		cloudProviderWithRegion := cloudProviderRegionModel{
			Id:        types.StringValue(cloudProvider.Id),
			Name:      types.StringValue(cloudProvider.Name),
			ShortName: types.StringValue(cloudProvider.ShortName),
		}

		list, diags := types.ListValueFrom(ctx, types.StringType, cloudProvider.Regions)
		if diags.Append(diags...); diags.HasError() {
			return
		}
		cloudProviderWithRegion.Regions = list
		state.CloudProviderRegions = append(state.CloudProviderRegions, cloudProviderWithRegion)
	}
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
