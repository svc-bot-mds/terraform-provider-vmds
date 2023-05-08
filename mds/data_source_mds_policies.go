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
)

var (
	_ datasource.DataSource              = &mdsPoliciesDatasource{}
	_ datasource.DataSourceWithConfigure = &mdsPoliciesDatasource{}
)

// instanceTypesDataSourceModel maps the data source schema data.
type mdsPoliciesDatasourceModel struct {
	Policies []mdsPoliciesModel `tfsdk:"policies"`
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
		Attributes: map[string]schema.Attribute{
			"names": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"policies": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
							Optional: true,
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
