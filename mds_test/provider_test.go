package mds_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/svc-bot-mds/terraform-provider-vmds/mds"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the VMDS client is properly configured.
	providerConfig = `
provider "vmds" {
   host     = "MDS_HOST_URL"
   api_token = "API_TOKEN"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"vmds": providerserver.NewProtocol6WithError(mds.New()),
	}
)
