package mds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_users" "users" {
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_users.users", "id"),
					resource.TestCheckResourceAttr("data.vmds_users.users", "users.#", "10"),
					resource.TestCheckResourceAttr("data.vmds_users.users", "users.0.email", "ptendolkar@vmware.com"),
					resource.TestCheckResourceAttr("data.vmds_users.users", "id", "placeholder"),
				),
			},
		},
	})
}
