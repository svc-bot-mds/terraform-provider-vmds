package mds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsRolesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_roles" "roles" {
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_roles.roles", "id"),
					resource.TestCheckResourceAttr("data.vmds_roles.roles", "roles.#", "5"),
					resource.TestCheckResourceAttr("data.vmds_roles.roles", "roles.0.name", "Operator"),
					resource.TestCheckResourceAttr("data.vmds_roles.roles", "roles.0.role_id", "StgManagedDataService:Operator"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.vmds_roles.roles", "id", "placeholder"),
				),
			},
		},
	})
}
