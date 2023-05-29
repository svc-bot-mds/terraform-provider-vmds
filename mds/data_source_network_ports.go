package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
)

var (
	_ datasource.DataSource              = &networkPortsDataSource{}
	_ datasource.DataSourceWithConfigure = &networkPortsDataSource{}
)

// networkPortsDatasource maps the datasource schema
type networkPortsDataSourceModel struct {
	NetworkPorts []networkPortsModel `tfsdk:"network_ports"`
	Id           types.String        `tfsdk:"id"`
}

type networkPortsModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Port        types.Int64  `tfsdk:"port"`
}

func NewNetworkPortsDataSource() datasource.DataSource {
	return &networkPortsDataSource{}
}

type networkPortsDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *networkPortsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_ports"
}

// Schema defines the schema for the data source.
func (d *networkPortsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"network_ports": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"port": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *networkPortsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networkPortsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	networkPorts, err := d.client.ServiceMetadata.GetNetworkPorts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS InstanceTypes",
			err.Error(),
		)
		return
	}

	// Map networkPorts body to model
	for _, networkPort := range networkPorts {
		networkPortsState := networkPortsModel{
			ID:          types.StringValue(networkPort.ID),
			Description: types.StringValue(networkPort.Description),
			Name:        types.StringValue(networkPort.Name),
			Port:        types.Int64Value(networkPort.Port),
		}

		state.NetworkPorts = append(state.NetworkPorts, networkPortsState)
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
func (d *networkPortsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
