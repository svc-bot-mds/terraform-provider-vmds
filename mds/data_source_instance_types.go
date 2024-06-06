package mds

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/service_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/controller"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"strconv"
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
	Id            types.String         `tfsdk:"id"`
	//Type          types.String         `tfsdk:"type"`
}

// instanceTypesModel maps coffees schema data.
type instanceTypesModel struct {
	ID           string                     `tfsdk:"id"`
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
		Description: "Used to fetch all instance sizes available for a service type on MDS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource.",
			},
			"service_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Type of the service. Supported values: %s .", supportedServiceTypesMarkdown()),
				Required:            true,
			},
			"instance_types": schema.ListNestedAttribute{
				Description: "List of the instance sizes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the size.",
							Computed:    true,
						},
						"instance_size": schema.StringAttribute{
							Description: "Friendly identifier of the size.",
							Computed:    true,
						},
						"instance_size_description": schema.StringAttribute{
							Description: "Description of the size.",
							Computed:    true,
						},
						"service_type": schema.StringAttribute{
							Description: "Type of the service supporting this size.",
							Required:    true,
						},
						"cpu": schema.StringAttribute{
							Description: "CPU that will be required by this size.",
							Computed:    true,
						},
						"memory": schema.StringAttribute{
							Description: "Memory that will be required by this size.",
							Computed:    true,
						},
						"storage": schema.StringAttribute{
							Description: "Storage that will be required by this size.",
							Computed:    true,
						},
						"metadata": schema.SingleNestedAttribute{
							Description: "Service specific additional resources.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"max_connections": schema.NumberAttribute{
									Description: "Total number of connections that can be established to the instance of this size.",
									Computed:    true,
								},
								"nodes": schema.NumberAttribute{
									Description: "Number of nodes that will be spawn of the instance of this size.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

//Read refreshes the Terraform state with the latest data.

func (d *instanceTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state *instanceTypesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if err := service_type.ValidateRoleType(state.ServiceType.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"invalid type",
			err.Error())
		return
	}
	query := &controller.MdsInstanceTypesQuery{
		ServiceType: state.ServiceType.ValueString(),
	}
	serviceInstanceTypes, err := d.client.Controller.GetServiceInstanceTypes(query)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TDH Service Types",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, instanceType := range serviceInstanceTypes.InstanceTypes {

		var nodes int64
		var maxconnections int64

		nodes, err := strconv.ParseInt(instanceType.Metadata.Nodes, 10, 64)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse Nodes value",
				err.Error(),
			)
			return
		}

		maxconnections, err = strconv.ParseInt(instanceType.Metadata.MaxConnections, 10, 64)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to parse MaxConnections value",
				err.Error(),
			)
			return
		}


		instanceTypesState := instanceTypesModel{
			ID:           instanceType.ID,
			InstanceSize: types.StringValue(instanceType.InstanceSize),
			Description:  types.StringValue(instanceType.SizeDescription),
			ServiceType:  types.StringValue(instanceType.ServiceType),
			CPU:          types.StringValue(instanceType.CPU),
			Memory:       types.StringValue(instanceType.Memory),
			Storage:      types.StringValue(instanceType.Storage),
			Metadata: instanceTypesMetadataModel{
				MaxConnections: types.Int64Value(maxconnections),
				Nodes:          types.Int64Value(nodes),
			},
		}

		state.InstanceTypes = append(state.InstanceTypes, instanceTypesState)
	}

	state.Id = types.StringValue(common.DataSource + common.InstanceTypesId)
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
