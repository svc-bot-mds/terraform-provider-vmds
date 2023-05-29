package mds_test

import (
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsNetworkPortsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_network_ports" "all" {
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_network_ports.all", "id"),
					resource.TestCheckResourceAttr("data.vmds_network_ports.all", "network_ports.#", "5"),
					resource.TestCheckResourceAttr("data.vmds_network_ports.all", "network_ports.0.name", "Metrics"),
					resource.TestCheckResourceAttr("data.vmds_network_ports.all", "network_ports.0.port", "443"),
					resource.TestCheckResourceAttr("data.vmds_network_ports.all", "id", common.DataSource+common.NetworkPortsId),
				),
			},
		},
	})
}
