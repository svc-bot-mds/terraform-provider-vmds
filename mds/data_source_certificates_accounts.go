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
	_ datasource.DataSource              = &certificatesDatasource{}
	_ datasource.DataSourceWithConfigure = &certificatesDatasource{}
)

// certificatesDatasourceModel maps the data source schema data.
type certificatesDatasourceModel struct {
	Id           types.String        `tfsdk:"id"`
	Certificates []certificatesModel `tfsdk:"certificates"`
}

type certificatesModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DomainName   types.String `tfsdk:"domain_name"`
	ProviderType types.String `tfsdk:"provider_type"`
	ExpiryTime   types.String `tfsdk:"expiry_time"`
}

// NewCertificatesDatasource is a helper function to simplify the provider implementation.
func NewCertificatesDatasource() datasource.DataSource {
	return &certificatesDatasource{}
}

// certificatesDatasource is the data source implementation.
type certificatesDatasource struct {
	client *mds.Client
}

// Metadata returns the data source type name.
func (d *certificatesDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

// Schema defines the schema for the data source.
func (d *certificatesDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Used to fetch all certificates on MDS for BYOC.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The testing framework requires an id attribute to be present in every data source and resource",
			},
			"certificates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the certificate.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the certificate.",
							Computed:    true,
						},
						"provider_type": schema.StringAttribute{
							Description: "Provider type of the certificate",
							Computed:    true,
						},
						"domain_name": schema.StringAttribute{
							Description: "Domain name of the certificate",
							Computed:    true,
						},
						"expiry_time": schema.StringAttribute{
							Description: "Expiry Time of the certificate",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *certificatesDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state certificatesDatasourceModel
	var certificateList []certificatesModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	query := &infra_connector.MDSCertificateQuery{}

	certificates, err := d.client.InfraConnector.GetCertificates(query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Byoc Certificates",
			err.Error(),
		)
		return
	}

	if certificates.Page.TotalPages > 1 {
		for i := 1; i <= certificates.Page.TotalPages; i++ {
			query.PageQuery.Index = i - 1
			totalCertificates, err := d.client.InfraConnector.GetCertificates(query)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Read Byoc certificates",
					err.Error(),
				)
				return
			}

			for _, certificateDto := range *totalCertificates.Get() {
				certificate := certificatesModel{
					ID:           types.StringValue(certificateDto.Id),
					Name:         types.StringValue(certificateDto.Name),
					DomainName:   types.StringValue(certificateDto.DomainName),
					ProviderType: types.StringValue(certificateDto.Provider),
					ExpiryTime:   types.StringValue(certificateDto.ExpiryTime),
				}
				certificateList = append(certificateList, certificate)
			}
		}

		tflog.Debug(ctx, "certificates dto", map[string]interface{}{"dto": certificateList})
		state.Certificates = append(state.Certificates, certificateList...)
	} else {
		for _, certificateDto := range *certificates.Get() {
			tflog.Info(ctx, "Converting certificate Dto1", map[string]interface{}{"dto": certificateDto})
			certificate := certificatesModel{
				ID:           types.StringValue(certificateDto.Id),
				Name:         types.StringValue(certificateDto.Name),
				DomainName:   types.StringValue(certificateDto.DomainName),
				ProviderType: types.StringValue(certificateDto.Provider),
				ExpiryTime:   types.StringValue(certificateDto.ExpiryTime),
			}
			tflog.Info(ctx, "converted certificate Dto", map[string]interface{}{"dto": certificate})
			state.Certificates = append(state.Certificates, certificate)
		}
	}
	state.Id = types.StringValue(common.DataSource + common.CertificateId)

	tflog.Info(ctx, "final", map[string]interface{}{"dto": state})
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *certificatesDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mds.Client)
}
