package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/policy_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	customer_metadata "github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
)

var (
	_ datasource.DataSource              = &mdsPoliciesDatasource{}
	_ datasource.DataSourceWithConfigure = &mdsPoliciesDatasource{}
)

// instanceTypesDataSourceModel maps the data source schema data.
type mdsPoliciesDatasourceModel struct {
	Policies []mdsPoliciesModel `tfsdk:"policies"`
	Id       types.String       `tfsdk:"id"`
	Names    types.List         `tfsdk:"names"`
}

// instanceTypesModel maps coffees schema data.
type mdsPoliciesModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// NewMdsPoliciesDatasource is a helper function to simplify the provider implementation.
func NewMdsPoliciesDatasource() datasource.DataSource {
	return &mdsPoliciesDatasource{}
}

// networkPoliciesDatasource is the data source implementation.
type mdsPoliciesDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *mdsPoliciesDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policies"
}

// Schema defines the schema for the data source.
func (d *mdsPoliciesDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Used to fetch all user access control policies for services.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource.",
			},
			"names": schema.ListAttribute{
				MarkdownDescription: "Names to search policies by. Ex: `[\"read-only-rmq\"]` .",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"policies": schema.ListNestedAttribute{
				Description: "List of fetched policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *mdsPoliciesDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mdsPoliciesDatasourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &customer_metadata.MdsPoliciesQuery{
		//TODO take policy type as input for type of service
		Type: policy_type.RABBITMQ,
	}

	nwPolicies, err := d.client.CustomerMetadata.GetPolicies(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read MDS Policies",
			err.Error(),
		)
		return
	}

	for _, mdsPolicyDTO := range *nwPolicies.Get() {
		policy := mdsPoliciesModel{
			ID:   types.StringValue(mdsPolicyDTO.ID),
			Name: types.StringValue(mdsPolicyDTO.Name),
		}
		tflog.Debug(ctx, "nwPolicy dto", map[string]interface{}{"dto": policy})
		state.Policies = append(state.Policies, policy)
	}

	state.Id = types.StringValue(common.DataSource + common.PoliciesId)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *mdsPoliciesDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
