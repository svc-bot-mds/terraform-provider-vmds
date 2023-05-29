package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/controller"
)

var (
	_ datasource.DataSource              = &clustersDatasource{}
	_ datasource.DataSourceWithConfigure = &clustersDatasource{}
)

// clustersDatasourceModel maps the data source schema data.
type clustersDatasourceModel struct {
	Clusters    []clustersModel `tfsdk:"clusters"`
	ID          types.String    `tfsdk:"id"`
	ServiceType types.String    `tfsdk:"service_type"`
}

// clustersModel maps clusters schema data.
type clustersModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// NewClustersDatasource is a helper function to simplify the provider implementation.
func NewClustersDatasource() datasource.DataSource {
	return &clustersDatasource{}
}

// clustersDatasource is the data source implementation.
type clustersDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *clustersDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

// Schema defines the schema for the data source.
func (d *clustersDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_type": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"clusters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *clustersDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clustersDatasourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &controller.MdsClustersQuery{
		ServiceType: state.ServiceType.ValueString(),
	}

	clusters, err := d.client.Controller.GetMdsClusters(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Clusters",
			err.Error(),
		)
		return
	}

	for _, mdsClusterDto := range *clusters.Get() {
		cluster := clustersModel{
			ID:   types.StringValue(mdsClusterDto.ID),
			Name: types.StringValue(mdsClusterDto.Name),
		}
		tflog.Debug(ctx, "mdsClusterDto dto", map[string]interface{}{"dto": cluster})
		state.Clusters = append(state.Clusters, cluster)
	}

	state.ID = types.StringValue("placeholder")
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *clustersDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
