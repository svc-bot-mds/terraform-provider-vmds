package mds

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/oauth_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/service_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"os"
	"strings"

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
	Host         types.String `tfsdk:"host"`
	Type         types.String `tfsdk:"type"`
	ApiToken     types.String `tfsdk:"api_token"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	OrgId        types.String `tfsdk:"org_id"`
}

// Metadata returns the provider type name.
func (p *mdsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vmds"
}

// Schema defines the provider-level schema for configuration data.
func (p *mdsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with VMware Managed Data Services",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for MDS API. May also be provided via MDS_HOST environment variable.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "OAuthType for the MDS API. It can be 'api_token' or 'client_credentials'",
				Required:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "API Token for MDS API. May also be provided via MDS_API_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client Id for MDS API. May also be provided via MDS_CLIENT_ID environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret for MDS API. May also be provided via MDS_CLIENT_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"org_id": schema.StringAttribute{
				Description: "Organization Id for MDS API. May also be provided via MDS_ORG_ID environment variable.",
				Optional:    true,
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

	//TODO read also from env variables
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown MDS API Host",
			"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_HOST environment variable.",
		)
	}
	if config.Type.ValueString() == oauth_type.ApiToken {
		if config.ApiToken.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ApiToken"),
				"Unknown MDS API Token",
				"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API Token. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_API_TOKEN environment variable.",
			)
		}
	}

	if config.Type.ValueString() == oauth_type.ClientCredentials {
		if config.ClientId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ClientId"),
				"Unknown MDS API ClientId",
				"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API ClientId. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_CLIENT_ID environment variable.",
			)
		}

		if config.ClientSecret.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_secret"),
				"Unknown MDS API Client Secret",
				"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API password. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_CLIENT_SECRET environment variable.",
			)
		}

		if config.OrgId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("org id"),
				"Unknown MDS API Org Id",
				"The provider cannot create the MDS API client as there is an unknown configuration value for the MDS API Org Id. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the MDS_ORG_ID environment variable.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("MDS_HOST")
	apiToken := os.Getenv("MDS_API_TOKEN")
	clientSecret := os.Getenv("MDS_CLIENT_SECRET")
	clientId := os.Getenv("MDS_CLIENT_ID")
	orgId := os.Getenv("MDS_ORG_ID")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if config.Type.ValueString() == oauth_type.ApiToken {
		if !config.ApiToken.IsNull() {
			apiToken = config.ApiToken.ValueString()
		}
	}
	if config.Type.ValueString() == oauth_type.ClientCredentials {
		if !config.ClientId.IsNull() {
			clientId = config.ClientId.ValueString()
		}

		if !config.ClientSecret.IsNull() {
			clientSecret = config.ClientSecret.ValueString()
		}

		if !config.OrgId.IsNull() {
			orgId = config.OrgId.ValueString()
		}
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

	if apiToken == "" && config.Type.ValueString() == oauth_type.ApiToken {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing MDS API Token",
			"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API Token. "+
				"Set the password value in the configuration or use the MDS_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if config.Type.ValueString() == oauth_type.ClientCredentials {
		if clientId == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_id"),
				"Missing MDS API Client Id",
				"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API Client Id. "+
					"Set the password value in the configuration or use the MDS_CLIENT_ID environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if clientSecret == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_secret"),
				"Missing MDS API Client Secret",
				"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API Client Secret. "+
					"Set the password value in the configuration or use the MDS_CLIENT_SECRET environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if orgId == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("org_id"),
				"Missing MDS API Org Id",
				"The provider cannot create the MDS API client as there is a missing or empty value for the MDS API Org Id. "+
					"Set the password value in the configuration or use the MDS_ORG_ID environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "mds_host", host)
	if config.Type.ValueString() == oauth_type.ClientCredentials {
		ctx = tflog.SetField(ctx, "mds_client_id", clientId)
		ctx = tflog.SetField(ctx, "mds_client_secret", clientSecret)
		ctx = tflog.SetField(ctx, "mds_org_id", orgId)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "mds_client_secret")
	} else {
		ctx = tflog.SetField(ctx, "mds_api_token", apiToken)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "mds_api_token")
	}

	tflog.Debug(ctx, "Creating MDS client")

	// Create a new MDS client using the configuration values
	client, err := mds.NewClient(&host, &model.ClientAuth{
		ApiToken:     apiToken,
		ClientSecret: clientSecret,
		ClientId:     clientId,
		OrgId:        orgId,
		OAuthAppType: config.Type.ValueString(),
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
		NewCloudAccountsDatasource,
		NewProviderTypesDataSource,
		NewCloudProviderRegionsDataSource,
		NewTshirtSizeDatasource,
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
		NewNetworkPolicyResource,
	}
}

func supportedServiceTypesMarkdown() string {
	var sb strings.Builder
	serviceTypes := service_type.GetAll()
	sb.WriteString(fmt.Sprintf("`%s`", serviceTypes[0]))
	for _, serviceType := range serviceTypes[1:] {
		sb.WriteString(fmt.Sprintf(", `%s`", serviceType))
	}
	return sb.String()
}
