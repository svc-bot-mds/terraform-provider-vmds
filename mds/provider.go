package mds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/oauth_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &mdsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &mdsProvider{}
}

// mdsProvider is the provider implementation.
type mdsProvider struct{}

// mdsProviderModel maps provider schema data to a Go type.
type mdsProviderModel struct {
	Host     types.String `tfsdk:"host"`
	ApiToken types.String `tfsdk:"api_token"`
}

// Metadata returns the provider type name.
func (p *mdsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vmds"
}

// Schema defines the provider-level schema for configuration data.
func (p *mdsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"api_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *mdsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MDS client")

	// Retrieve provider data from configuration
	var config mdsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown MDS API Host",
			"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_HOST environment variable.",
		)
	}

	if config.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown MDS API Password",
			"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("MDS_HOST")
	apiToken := os.Getenv("MDS_API_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiToken.IsNull() {
		apiToken = config.ApiToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing MDS API Host",
			"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API host. "+
				"Set the host value in the configuration or use the MDS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing MDS API Token",
			"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API Token. "+
				"Set the password value in the configuration or use the MDS_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "mds_host", host)
	ctx = tflog.SetField(ctx, "mds_api_token", apiToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "mds_api_token")

	tflog.Debug(ctx, "Creating MDS client")

	// Create a new MDS client using the configuration values
	client, err := mds.NewClient(&host, &model.ClientAuth{
		ApiToken:     apiToken,
		OAuthAppType: oauth_type.ApiToken,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create MDS API Client",
			"An unexpected error occurred when creating the MDS API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"MDS Client Error: "+err.Error(),
		)
		return
	}

	// Make the MDS client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured MDS client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *mdsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewInstanceTypesDataSource,
		NewRegionsDataSource,
		NewNetworkPoliciesDataSource,
		NewNetworkPortsDataSource,
		NewUsersDataSource,
		NewRolesDataSource,
		NewMdsPoliciesDatasource,
		NewServiceAccountsDataSource,
		NewPolicyTypesDataSource,
		NewClusterMetadataDataSource,
		NewClustersDatasource,
		NewServiceRolesDatasource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *mdsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewClusterNetworkPoliciesAssociationResource,
		NewUserResource,
		NewServiceAccountResource,
		NewPolicyResource,
	}
}
