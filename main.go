package main

import (
	"context"
	"github.com/svc-bot-mds/terraform-provider-vmds/mds"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), mds.New, providerserver.ServeOpts{
		Address: "svc-bot-mds/vmds",
	})
	if err != nil {
		return
	}
}
