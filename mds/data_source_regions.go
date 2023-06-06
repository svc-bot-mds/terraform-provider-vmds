package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/service_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/controller"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"strconv"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &regionsDataSource{}
	_ datasource.DataSourceWithConfigure = &regionsDataSource{}
)

type regionsDataSourceModel struct {
	Cpu                types.String   `tfsdk:"cpu"`
	Provider           types.String   `tfsdk:"cloud_provider"`
	Memory             types.String   `tfsdk:"memory"`
	Storage            types.String   `tfsdk:"storage"`
	NodeCount          types.String   `tfsdk:"node_count"`
	InstanceSize       types.String   `tfsdk:"instance_size"`
	Regions            []RegionsModel `tfsdk:"regions"`
	DedicatedDataPlane types.Bool     `tfsdk:"dedicated_data_plane"`
	Id                 types.String   `tfsdk:"id"`
}

type RegionsModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DataPlaneIds types.List   `tfsdk:"data_plane_ids"`
}

func NewRegionsDataSource() datasource.DataSource {
	return &regionsDataSource{}
}

type regionsDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *regionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (d *regionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch the regions having data-planes by desired amount of resources available.\n" +
			"## Note:\n" +
			"- At a time, either `instance_size` or all of (`cpu`, `memory`, `storage`, `node_count`) can be passed.",
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "Shortname of cloud provider platform where data-plane lives. Ex: `aws`, `gcp` .",
				Required:            true,
			},
			"cpu": schema.StringAttribute{
				MarkdownDescription: "K8s CPU units required. Ex: `500m`, `1` (1000m) .",
				Optional:            true,
			},
			"memory": schema.StringAttribute{
				MarkdownDescription: "K8s memory units required. Ex: `800Mi`, `2Gi` .",
				Optional:            true,
			},
			"storage": schema.StringAttribute{
				MarkdownDescription: "K8s storage units required. Ex: `2Gi` .",
				Optional:            true,
			},
			"node_count": schema.StringAttribute{
				MarkdownDescription: "Count of worker nodes that must be present in a data-plane. Ex: `3` .",
				Optional:            true,
			},
			"instance_size": schema.StringAttribute{
				MarkdownDescription: "Type of instance size. Supported values: `XX-SMALL`, `X-SMALL`, `SMALL`, `LARGE`, `XX-LARGE`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("cpu"),
						path.MatchRoot("memory"),
						path.MatchRoot("storage"),
						path.MatchRoot("node_count"),
					}...),
				},
			},
			"dedicated_data_plane": schema.BoolAttribute{
				MarkdownDescription: "If set to `true`, only data-planes that are exclusive to current Org (determined by API token used) are queried. Else only shared ones.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource.",
			},
			"regions": schema.ListNestedAttribute{
				Description: "Response of regional data-planes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the region.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the region.",
							Computed:    true,
						},
						"data_plane_ids": schema.ListAttribute{
							Description: "List of data-plane IDs that are created by SRE & registered on MDS.",
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
func (d *regionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *regionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state regionsDataSourceModel

	//Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	regionQuery := &infra_connector.DataPlaneRegionsQuery{
		CPU:       state.Cpu.ValueString(),
		NodeCount: state.NodeCount.ValueString(),
		Memory:    state.Memory.ValueString(),
		Storage:   state.Storage.ValueString(),
		Provider:  state.Provider.ValueString(),
	}
	var typeDetail model.MdsInstanceType
	if !state.InstanceSize.IsNull() {
		instanceTypes, err := d.client.Controller.GetServiceInstanceTypes(&controller.MdsInstanceTypesQuery{
			ServiceType: service_type.RABBITMQ,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read MDS Regions:",
				err.Error(),
			)
			return
		}
		for _, instanceType := range instanceTypes.InstanceTypes {
			if state.InstanceSize.ValueString() == instanceType.InstanceSize {
				typeDetail = instanceType
				break
			}
		}
	}
	regionQuery.CPU = typeDetail.CPU
	regionQuery.Memory = typeDetail.Memory
	regionQuery.Storage = typeDetail.Storage
	regionQuery.NodeCount = strconv.Itoa(int(typeDetail.Metadata.Nodes))
	if state.DedicatedDataPlane.ValueBool() {
		regionQuery.OrgId = d.client.Root.OrgId
	}
	regions, err := d.client.InfraConnector.GetRegionsWithDataPlanes(regionQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Regions:",
			err.Error(),
		)
		return
	}
	if saveRegionsToState(&ctx, &resp.Diagnostics, &state, regions) != 0 {
		return
	}
	state.Id = types.StringValue(common.DataSource + common.RegionsId)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func saveRegionsToState(ctx *context.Context, diagnostics *diag.Diagnostics, state *regionsDataSourceModel, regions map[string][]string) int8 {
	for regionName, dataPlaneIds := range regions {
		instanceTypesState := RegionsModel{
			ID:   types.StringValue(regionName),
			Name: types.StringValue(regionName),
		}
		list, diags := types.ListValueFrom(*ctx, types.StringType, dataPlaneIds)
		if diagnostics.Append(diags...); diagnostics.HasError() {
			return 1
		}
		instanceTypesState.DataPlaneIds = list
		state.Regions = append(state.Regions, instanceTypesState)
	}
	return 0
}
