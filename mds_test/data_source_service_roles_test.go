package mds_test

import (
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsServiceRolesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_service_roles" "roles"{
  type = "RABBITMQ"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_service_roles.roles", "id"),
					resource.TestCheckResourceAttrSet("data.vmds_service_roles.roles", "type"),
					resource.TestCheckResourceAttr("data.vmds_service_roles.roles", "roles.#", "6"),
					resource.TestCheckResourceAttr("data.vmds_service_roles.roles", "roles.0.name", "write"),
					resource.TestCheckResourceAttr("data.vmds_service_roles.roles", "roles.0.role_id", "StgManagedDataService:RMQWrite"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.vmds_service_roles.roles", "id", common.DataSource+common.ServiceRolesId),
				),
			},
		},
	})
}
