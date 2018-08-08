package opentelekomcloud

import (
"fmt"
"testing"

"github.com/hashicorp/terraform/helper/resource"
"github.com/hashicorp/terraform/terraform"
"github.com/huaweicloud/golangsdk/openstack/cce/v1/clusters"
)

func TestAccOTCCCEClusterV1_basic(t *testing.T) {
	var cluster clusters.RetrievedCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOTCCheckCCEClusterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCCCEClusterV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCCheckCCEClusterV1Exists("opentelekomcloud_cce_cluster_v1.cluster_1", &cluster),
				),
			},
			resource.TestStep{
				Config: testAccOTCCCEClusterV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCCheckCCEClusterV1Exists("opentelekomcloud_cce_cluster_v1.cluster_1", &cluster),
				),
			},
		},
	})
}

func TestAccOTCCCEClusterV1_timeout(t *testing.T) {
	var cluster clusters.RetrievedCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOTCCheckCCEClusterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCCCEClusterV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCCheckCCEClusterV1Exists("opentelekomcloud_cce_cluster_v1.cluster_1", &cluster),
				),
			},
		},
	})
}

func testAccOTCCheckCCEClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccOTCCheckCCEClusterV1Exists(n string, cluster *clusters.RetrievedCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.ID != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

const testAccOTCCCEClusterV1_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"

}
resource "opentelekomcloud_cce_cluster_v1" "cluster_1" {
  name = "opentelekomcloud-cluster"
  vpc_id="${opentelekomcloud_vpc_v1.vpc_1.id}"
  subnet_id="${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  cluster_type="Single"
  description="test cluster"
 }
`

const testAccOTCCCEClusterV1_update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"

}
resource "opentelekomcloud_cce_cluster_v1" "cluster_1" {
  name = "opentelekomcloud-cluster"
  vpc_id="${opentelekomcloud_vpc_v1.vpc_1.id}"
  subnet_id="${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  cluster_type="Single"
  description="New description"
  }
`
const testAccOTCCCEClusterV1_timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"

}
resource "opentelekomcloud_cce_cluster_v1" "cluster_1" {
  name = "opentelekomcloud-cluster"
  vpc_id="${opentelekomcloud_vpc_v1.vpc_1.id}"
  subnet_id="${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
  cluster_type="Single"
  description="test cluster"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

