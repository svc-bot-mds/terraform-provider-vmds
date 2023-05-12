package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
)

var (
	_ datasource.DataSource              = &clusterMetdataDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterMetdataDataSource{}
)

type clusterMetdataDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Name         types.String `tfsdk:"name"`
	ServiceType  types.String `tfsdk:"service_type"`
	Status       types.String `tfsdk:"status"`
}

func NewClusterMetadataDataSource() datasource.DataSource {
	return &clusterMetdataDataSource{}
}

type clusterMetdataDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *clusterMetdataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_metadata"
}

// Schema defines the schema for the data source.
func (d *clusterMetdataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"provider_name": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"service_type": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *clusterMetdataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clusterMetdataDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	clusterMetadata, err := d.client.Controller.GetMdsClusterMetaData(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Cluster Metadata",
			err.Error(),
		)
		return
	}

	// Map cluster metadata body to model
	metadataDetails := clusterMetdataDataSourceModel{
		Id:           types.StringValue(clusterMetadata.Id),
		Name:         types.StringValue(clusterMetadata.Name),
		ProviderName: types.StringValue(clusterMetadata.Provider),
		ServiceType:  types.StringValue(clusterMetadata.ServiceType),
		Status:       types.StringValue(clusterMetadata.Status),
	}

	state = metadataDetails
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *clusterMetdataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
