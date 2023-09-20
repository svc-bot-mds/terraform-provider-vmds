package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

var (
	_ datasource.DataSource              = &tshirtSizeDatasource{}
	_ datasource.DataSourceWithConfigure = &tshirtSizeDatasource{}
)

// tshirtSizeDatasourceModel maps the data source schema data.
type tshirtSizeDatasourceModel struct {
	Id         types.String      `tfsdk:"id"`
	TshirtSize []tshirtSizeModel `tfsdk:"tshirt_sizes"`
}

type tshirtSizeModel struct {
	Name     types.String `tfsdk:"name"`
	Nodes    types.Int64  `tfsdk:"nodes"`
	Provider types.String `tfsdk:"provider"`
	Storage  types.String `tfsdk:"storage"`
	Type     types.String `tfsdk:"type"`
}

// NewCloudAccountsDatasource is a helper function to simplify the provider implementation.

func NewTshirtSizeDatasource() datasource.DataSource {
	return &tshirtSizeDatasource{}
}

// tshirtSizeDatasource is the data source implementation.
type tshirtSizeDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *tshirtSizeDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tshirt_size"
}

// Schema defines the schema for the data source.
func (d *tshirtSizeDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch all tshirt size available for BYOC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"tshirt_sizes": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"storage": schema.StringAttribute{
							Description: "storage of the tshirt.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the tshirt.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the tshirt",
							Computed:    true,
						},
						"nodes": schema.Int64Attribute{
							Description: "Nodes available for the tshirt",
							Computed:    true,
						},
						"provider": schema.StringAttribute{
							Description: "Provider for the tshirt",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tshirtSizeDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tshirtSizeDatasourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &infra_connector.MdsTshirtSizesQuery{}

	tshirtSizes, err := d.client.InfraConnector.GetTshirtSizes(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read BYOC Tshirt sizes",
			err.Error(),
		)
		return
	}

	for _, cloudAccountDto := range *tshirtSizes.Get() {
		tflog.Info(ctx, "Converting tshirt size dto")
		tshirt := tshirtSizeModel{
			Nodes:    types.Int64Value(cloudAccountDto.Nodes),
			Name:     types.StringValue(cloudAccountDto.Name),
			Provider: types.StringValue(cloudAccountDto.Provider),
			Storage:  types.StringValue(cloudAccountDto.Storage),
			Type:     types.StringValue(cloudAccountDto.Type),
		}
		tflog.Debug(ctx, "converted tshirt size dto", map[string]interface{}{"dto": tshirt})
		state.TshirtSize = append(state.TshirtSize, tshirt)
	}
	state.Id = types.StringValue(common.DataSource + common.CloudAccountsId)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *tshirtSizeDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
