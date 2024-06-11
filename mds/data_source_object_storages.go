package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	infra_connector "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

var (
	_ datasource.DataSource              = &objectStorageDatasource{}
	_ datasource.DataSourceWithConfigure = &objectStorageDatasource{}
)

// objectStoragesDatasourceModel maps the data source schema data.
type objectStoragesDatasourceModel struct {
	Id             types.String         `tfsdk:"id"`
	ObjectStorages []objectStorageModel `tfsdk:"list"`
}

type objectStorageModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	BucketName types.String `tfsdk:"bucket_name"`
	Endpoint   types.String `tfsdk:"endpoint"`
	Region     types.String `tfsdk:"region"`
}

// NewObjectStorageDatasource is a helper function to simplify the provider implementation.
func NewObjectStorageDatasource() datasource.DataSource {
	return &objectStorageDatasource{}
}

// objectStorageDatasource is the data source implementation.
type objectStorageDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *objectStorageDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storages"
}

// Schema defines the schema for the data source.
func (d *objectStorageDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch all object storages on MDS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the object storage.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the object storage.",
							Computed:    true,
						},
						"endpoint": schema.StringAttribute{
							Description: "Endpoint of the object storage.",
							Computed:    true,
						},
						"bucket_name": schema.StringAttribute{
							Description: "Name of the initial bucket to create.",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region where object storage is created.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *objectStorageDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectStoragesDatasourceModel
	var modelList []objectStorageModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &infra_connector.ObjectStorageQuery{}

	response, err := d.client.InfraConnector.GetObjectStorages(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Object Storages",
			err.Error(),
		)
		return
	}

	if response.Page.TotalPages > 1 {
		for i := 1; i <= response.Page.TotalPages; i++ {
			query.PageQuery.Index = i - 1
			page, err := d.client.InfraConnector.GetObjectStorages(query)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Read Object Storages",
					err.Error(),
				)
				return
			}

			for _, dto := range *page.Get() {
				tfModel := d.convertToTfModel(dto)
				modelList = append(modelList, tfModel)
			}
		}

		tflog.Debug(ctx, "object storage models: ", map[string]interface{}{"models": modelList})
		state.ObjectStorages = append(state.ObjectStorages, modelList...)
	} else {
		for _, dto := range *response.Get() {
			tflog.Info(ctx, "Converting dto: ", map[string]interface{}{"dto": dto})
			tfModel := d.convertToTfModel(dto)
			tflog.Info(ctx, "converted object storage model: ", map[string]interface{}{"model": tfModel})
			state.ObjectStorages = append(state.ObjectStorages, tfModel)
		}
	}
	state.Id = types.StringValue(common.DataSource + common.ObjectStorageId)

	tflog.Info(ctx, "final", map[string]interface{}{"dto": state})
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *objectStorageDatasource) convertToTfModel(dto model.MdsObjectStorage) objectStorageModel {
	return objectStorageModel{
		ID:         types.StringValue(dto.Id),
		Name:       types.StringValue(dto.Name),
		Endpoint:   types.StringValue(dto.Endpoint),
		BucketName: types.StringValue(dto.BucketName),
		Region:     types.StringValue(dto.Region),
	}
}

// Configure adds the provider configured client to the data source.
func (d *objectStorageDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
