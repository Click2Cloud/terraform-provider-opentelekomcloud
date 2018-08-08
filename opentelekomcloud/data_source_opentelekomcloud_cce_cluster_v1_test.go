package opentelekomcloud

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)


func TestAccOpenTelekomCloudCCEClusterV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudRtsStackV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsStackV1DataSourceID("data.opentelekomcloud_cce_cluster_v1.clusters"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "name", "opentelekomcloud-cce"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "status", "EMPTY"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v1.clusters", "type", "Single"),
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
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "opentelekomcloud-vpc"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud-subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.20.1"
  dhcp_enable = "true"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
}

resource "opentelekomcloud_cce_cluster_v1" "cluster_1" {
  kind = "cluster"
  name = "opentelekomcloud-cce"
  api_version = "v1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  subnet_id = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  type = "Single"
}

data "opentelekomcloud_cce_cluster_v1" "clusters" {
  name = "${opentelekomcloud_cce_cluster_v1.cluster_1.name}"
}
`
