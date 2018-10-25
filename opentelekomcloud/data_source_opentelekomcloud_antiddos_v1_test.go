package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAntiDdosV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAntiDdosV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1DataSourceID("data.opentelekomcloud_antiddos_v1.antiddos"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_antiddos_v1.antiddos", "floating_ip_id", "160.44.206.31"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_antiddos_v1.antiddos", "status", "normal"),
				),
			},
		},
	})
}

func testAccCheckAntiDdosV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

const testAccAntiDdosV1DataSource_basic = `
data "opentelekomcloud_antiddos_v1" "antiddos" {  
  id = "7deb25d6-7a56-4d0b-9f28-305a069037d2"
}
`
