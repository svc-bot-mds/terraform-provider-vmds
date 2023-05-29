package mds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsClustersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_clusters" "cluster_list"{
  											service_type = "RABBITMQ"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_clusters.cluster_list", "id"),
					resource.TestCheckResourceAttrSet("data.vmds_clusters.cluster_list", "service_type"),
					resource.TestCheckResourceAttr("data.vmds_clusters.cluster_list", "clusters.#", "26"),
					resource.TestCheckResourceAttr("data.vmds_clusters.cluster_list", "clusters.0.name", "audit-test-dnd"),
					resource.TestCheckResourceAttr("data.vmds_clusters.cluster_list", "id", "placeholder"),
				),
			},
		},
	})
}
