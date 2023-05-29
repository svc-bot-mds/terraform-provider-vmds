package mds_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* Set up any prerequisites or check for required dependencies */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					locals {
						service_type  = "RABBITMQ"
						provider      = "aws"
						instance_type      = "XX-SMALL"
						region             = "eu-west-1"
					}

					resource "vmds_cluster" "test" {
						name               = "my-rmq-cls-test-terraform5"
						service_type       = local.service_type
						cloud_provider     = local.provider
						instance_size      = local.instance_type
						region             = local.region
						network_policy_ids = ["646f030f8c626b5a2b59d158"]
						tags               = ["mds-tf", "example", "new-tag", "create"]
						timeouts = {
							create = "1m"
							delete = "1m"
						}

						lifecycle {
							ignore_changes = [instance_size, name, cloud_provider, region, service_type]
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Add validation checks here to verify the resource state
					resource.TestCheckResourceAttr("vmds_cluster.test", "name", "my-rmq-cls-test-terraform5"),
					resource.TestCheckResourceAttr("vmds_cluster.test", "service_type", "RABBITMQ"),
				),
			},
			{
				Config: providerConfig + `
			locals {
						service_type  = "RABBITMQ"
						provider      = "aws"
						policy_with_create = ["open-to-all"]
						instance_type      = "XX-SMALL"
						region             = "eu-west-1"
					}

					resource "vmds_cluster" "test" {
						name               = "my-rmq-cls-test-terraform5"
						service_type       = local.service_type
						cloud_provider     = local.provider
						instance_size      = local.instance_type
						region             = local.region
						network_policy_ids = ["646f030f8c626b5a2b59d158"]
						tags               = ["mds-tf", "example", "new-tag", "edit"]
						timeouts = {
							create = "1m"
							delete = "1m"
						}

						lifecycle {
							ignore_changes = [instance_size, name, cloud_provider, region, service_type]
						}
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vmds_cluster.test", "tags.#", "4"),
					resource.TestCheckResourceAttr("vmds_cluster.test", "tags.0", "edit"),
				),
			},

			{
				ResourceName:      "vmds_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						cluster := s.RootModule().Resources["vmds_cluster.test"]
						status := cluster.Primary.Attributes["status"]

						// Perform status validation checks
						if status != "DELETE_IN_PROGRESS" && status != "DELETED" {
							return fmt.Errorf("unexpected status. Expected: %s, Or %s, Got: %s", "DELETE_IN_PROGRESS", "DELETED", status)
						}

						return nil
					},
				),
			},
		},
	})
}
