package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &policyTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &policyTypesDataSource{}
)

type policyTypesDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	PolicyTypes []string     `tfsdk:"policy_types"`
}

func NewPolicyTypesDataSource() datasource.DataSource {
	return &policyTypesDataSource{}
}

type policyTypesDataSource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *policyTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_types"
}

func (d *policyTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"policy_types": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *policyTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *policyTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "getPolicyTypes")
	var state policyTypesDataSourceModel

	//Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	tflog.Info(ctx, "getPolicyTypes")
	typesList, err := d.client.ServiceMetadata.GetPolicyTypes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Policy Types:",
			err.Error(),
		)
		return
	}
	state.Id = types.StringValue("placeholder")
	state.PolicyTypes = append(state.PolicyTypes, typesList...)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
