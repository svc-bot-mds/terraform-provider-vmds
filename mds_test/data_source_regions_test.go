package mds_test

import (
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsRegionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_regions" "available_regions" {
  cpu                  = "1"
  cloud_provider       = "aws"
  memory               = "4Gi"
  storage              = "4Gi"
  node_count           = "1"
  dedicated_data_plane = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "id"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "cloud_provider"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "cpu"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "dedicated_data_plane"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "memory"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "node_count"),
					resource.TestCheckResourceAttrSet("data.vmds_regions.available_regions", "storage"),
					resource.TestCheckResourceAttr("data.vmds_regions.available_regions", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.vmds_regions.available_regions", "regions.0.name", "eu-west-1"),
					resource.TestCheckResourceAttr("data.vmds_regions.available_regions", "id", common.DataSource+common.RegionsId),
				),
			},
		},
	})
}
