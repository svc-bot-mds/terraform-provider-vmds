package mds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsClusterMetadataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_cluster_metadata" "metadata" {
    										id = "dummyid"
  				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmds_cluster_metadata.metadata", "id", "dummmyid"),
					resource.TestCheckResourceAttr("data.vmds_cluster_metadata.metadata", "name", "test"),
					resource.TestCheckResourceAttr("data.vmds_cluster_metadata.metadata", "provider_name", "aws"),
					resource.TestCheckResourceAttr("data.vmds_cluster_metadata.metadata", "service_type", "RABBITMQ"),
					resource.TestCheckResourceAttr("data.vmds_cluster_metadata.metadata", "status", "READY"),
				),
			},
		},
	})
}
