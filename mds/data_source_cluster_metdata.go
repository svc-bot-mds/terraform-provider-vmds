package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
)

var (
	_ datasource.DataSource              = &clusterMetadataDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterMetadataDataSource{}
)

type clusterMetadataDataSourceModel struct {
	Id           types.String     `tfsdk:"id"`
	ProviderName types.String     `tfsdk:"provider_name"`
	Name         types.String     `tfsdk:"name"`
	ServiceType  types.String     `tfsdk:"service_type"`
	Status       types.String     `tfsdk:"status"`
	Vhosts       []VHostsModel    `tfsdk:"vhosts"`
	Queues       []QueuesModel    `tfsdk:"queues"`
	Exchanges    []ExchangesModel `tfsdk:"exchanges"`
	Bindings     []BindingsModel  `tfsdk:"bindings"`
}

type VHostsModel struct {
	Name types.String `tfsdk:"name"`
}

type QueuesModel struct {
	Name  types.String `tfsdk:"name"`
	VHost types.String `tfsdk:"vhost"`
}

type ExchangesModel struct {
	Name  types.String `tfsdk:"name"`
	VHost types.String `tfsdk:"vhost"`
}

type BindingsModel struct {
	Source          types.String `tfsdk:"source"`
	VHost           types.String `tfsdk:"vhost"`
	RoutingKey      types.String `tfsdk:"routing_key"`
	Destination     types.String `tfsdk:"destination"`
	DestinationType types.String `tfsdk:"destination_type"`
}

func NewClusterMetadataDataSource() datasource.DataSource {
	return &clusterMetadataDataSource{}
}

type clusterMetadataDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *clusterMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_metadata"
}

// Schema defines the schema for the data source.
func (d *clusterMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch metadata of a cluster by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the cluster.",
				Required:    true,
			},
			"provider_name": schema.StringAttribute{
				Description: "Name of the data-plane's cloud provider where cluster is deployed.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the cluster.",
				Computed:    true,
			},
			"service_type": schema.StringAttribute{
				Description: "Type of the service of the cluster.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the cluster.",
				Computed:    true,
			},
			"vhosts": schema.ListNestedAttribute{
				Description: "List of the vHosts. Specific to `RABBITMQ` service.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the vHost.",
							Computed:    true,
						},
					},
				},
			},
			"queues": schema.ListNestedAttribute{
				Description: "List of the Queues. Specific to `RABBITMQ` service.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the queue.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "vHost of the queue.",
							Computed:    true,
						},
					},
				},
			},
			"exchanges": schema.ListNestedAttribute{
				Description: "List of the Exchanges. Specific to `RABBITMQ` service.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the exchange.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "vHost of the exchange.",
							Computed:    true,
						},
					},
				},
			},
			"bindings": schema.ListNestedAttribute{
				Description: "List of the Bindings. Specific to `RABBITMQ` service.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": schema.StringAttribute{
							Description: "Source exchange.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "vHost name.",
							Computed:    true,
						},
						"routing_key": schema.StringAttribute{
							Description: "Routing key.",
							Computed:    true,
						},
						"destination": schema.StringAttribute{
							Description: "Destination exchange.",
							Computed:    true,
						},
						"destination_type": schema.StringAttribute{
							Description: "Type of the destination.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *clusterMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clusterMetadataDataSourceModel

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
	metadataDetails := clusterMetadataDataSourceModel{
		Id:           types.StringValue(clusterMetadata.Id),
		Name:         types.StringValue(clusterMetadata.Name),
		ProviderName: types.StringValue(clusterMetadata.Provider),
		ServiceType:  types.StringValue(clusterMetadata.ServiceType),
		Status:       types.StringValue(clusterMetadata.Status),
	}

	if len(clusterMetadata.VHosts) > 0 {
		// Map response body to model
		for _, vhost := range clusterMetadata.VHosts {
			vhosts := VHostsModel{
				Name: types.StringValue(vhost.Name),
			}
			metadataDetails.Vhosts = append(state.Vhosts, vhosts)
		}
	}
	if len(clusterMetadata.Bindings) > 0 {
		for _, binding := range clusterMetadata.Bindings {
			bindings := BindingsModel{
				Source:          types.StringValue(binding.Source),
				DestinationType: types.StringValue(binding.DestinationType),
				Destination:     types.StringValue(binding.Destination),
				VHost:           types.StringValue(binding.VHost),
				RoutingKey:      types.StringValue(binding.RoutingKey),
			}
			metadataDetails.Bindings = append(state.Bindings, bindings)
		}
	}
	if len(clusterMetadata.Queues) > 0 {
		for _, queue := range clusterMetadata.Queues {
			queues := QueuesModel{
				Name:  types.StringValue(queue.Name),
				VHost: types.StringValue(queue.VHost),
			}
			metadataDetails.Queues = append(state.Queues, queues)
		}
	}
	if len(clusterMetadata.Exchanges) > 0 {
		for _, exchange := range clusterMetadata.Exchanges {
			exchanges := ExchangesModel{
				Name:  types.StringValue(exchange.Name),
				VHost: types.StringValue(exchange.VHost),
			}
			metadataDetails.Exchanges = append(state.Exchanges, exchanges)
		}
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
func (d *clusterMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
