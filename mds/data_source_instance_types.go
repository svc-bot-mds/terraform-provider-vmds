package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/service_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/controller"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &instanceTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &instanceTypesDataSource{}
)

// instanceTypesDataSourceModel maps the data source schema data.
type instanceTypesDataSourceModel struct {
	InstanceTypes []instanceTypesModel `tfsdk:"instance_types"`
	ServiceType   types.String         `tfsdk:"service_type"`
}

// instanceTypesModel maps coffees schema data.
type instanceTypesModel struct {
	ID           types.String               `tfsdk:"id"`
	InstanceSize types.String               `tfsdk:"instance_size"`
	Description  types.String               `tfsdk:"instance_size_description"`
	ServiceType  types.String               `tfsdk:"service_type"`
	CPU          types.String               `tfsdk:"cpu"`
	Memory       types.String               `tfsdk:"memory"`
	Storage      types.String               `tfsdk:"storage"`
	Metadata     instanceTypesMetadataModel `tfsdk:"metadata"`
}

// instanceTypesMetadataModel maps instanceType metadata data
type instanceTypesMetadataModel struct {
	MaxConnections types.Int64 `tfsdk:"max_connections"`
	Nodes          types.Int64 `tfsdk:"nodes"`
}

// NewInstanceTypesDataSource is a helper function to simplify the provider implementation.
func NewInstanceTypesDataSource() datasource.DataSource {
	return &instanceTypesDataSource{}
}

// instanceTypesDataSource is the data source implementation.
type instanceTypesDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *instanceTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_types"
}

// Schema defines the schema for the data source.
func (d *instanceTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_type": schema.StringAttribute{
				Required: true,
			},
			"instance_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"instance_size": schema.StringAttribute{
							Computed: true,
						},
						"instance_size_description": schema.StringAttribute{
							Computed: true,
						},
						"service_type": schema.StringAttribute{
							Required: true,
						},
						"cpu": schema.StringAttribute{
							Computed: true,
						},
						"memory": schema.StringAttribute{
							Computed: true,
						},
						"storage": schema.StringAttribute{
							Computed: true,
						},
						"metadata": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"max_connections": schema.NumberAttribute{
									Computed: true,
								},
								"nodes": schema.NumberAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *instanceTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instanceTypesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	serviceInstanceTypes, err := d.client.Controller.GetServiceInstanceTypes(&controller.MdsInstanceTypesQuery{
		ServiceType: service_type.RABBITMQ,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS InstanceTypes",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, instanceType := range serviceInstanceTypes.InstanceTypes {
		instanceTypesState := instanceTypesModel{
			ID:           types.StringValue(instanceType.ID),
			InstanceSize: types.StringValue(instanceType.InstanceSize),
			Description:  types.StringValue(instanceType.SizeDescription),
			ServiceType:  types.StringValue(instanceType.ServiceType),
			CPU:          types.StringValue(instanceType.CPU),
			Memory:       types.StringValue(instanceType.Memory),
			Storage:      types.StringValue(instanceType.Storage),
			Metadata: instanceTypesMetadataModel{
				MaxConnections: types.Int64Value(instanceType.Metadata.MaxConnections),
				Nodes:          types.Int64Value(instanceType.Metadata.Nodes),
			},
		}

		state.InstanceTypes = append(state.InstanceTypes, instanceTypesState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *instanceTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
