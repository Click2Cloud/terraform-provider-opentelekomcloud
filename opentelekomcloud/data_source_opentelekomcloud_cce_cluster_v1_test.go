package opentelekomcloud

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOpenTelekomCloudCCEClusterV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudRtsStackV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsStackV1DataSourceID("data.opentelekomcloud_cce_cluster_v1.clusters"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "name", "disha-test"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "status", "AVAILABLE"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "vpc_name", "vpc-disha"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "subnet", "subnet-28a5"),
				),
			},
		},
	})
}

func testAccCheckRtsStackV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

var testAccOpenTelekomCloudRtsStackV1DataSource_basic = `
resource "opentelekomcloud_cce_cluster_v1" "cluster_1" {
  
}
data "opentelekomcloud_cce_cluster_v1" "clusters" {
        name = "${opentelekomcloud_cce_cluster_v1.cluster_1.name}"
}
`
