package mds_test

import (
	"github.com/svc-bot-mds/terraform-provider-vmds/constants/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMdsInstanceTypesDataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "vmds_instance_types" "rmq" {
  service_type = "RABBITMQ"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmds_instance_types.rmq", "instance_types.#", "5"),
					resource.TestCheckResourceAttr("data.vmds_instance_types.rmq", "instance_types.0.instance_size", "LARGE"),
					resource.TestCheckResourceAttr("data.vmds_instance_types.rmq", "instance_types.0.service_type", "RABBITMQ"),
					resource.TestCheckResourceAttr("data.vmds_instance_types.rmq", "id", common.DataSource+common.InstanceTypesId),
				),
			},
		},
	})
}
