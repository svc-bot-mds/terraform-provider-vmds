package mds_test

import (
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_policies" "policies" {
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_policies.policies", "id"),
					resource.TestCheckResourceAttr("data.vmds_policies.policies", "policies.#", "27"),
					resource.TestCheckResourceAttr("data.vmds_policies.policies", "policies.0.name", "test-tfddwqe"),
					resource.TestCheckResourceAttr("data.vmds_policies.policies", "id", common.DataSource+common.PoliciesId),
				),
			},
		},
	})
}
