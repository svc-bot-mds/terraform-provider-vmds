package main

import (
	"context"
	"github.com/svc-bot-mds/terraform-provider-vmds/mds"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), mds.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/vmds",
	})
	if err != nil {
		return
	}
}
