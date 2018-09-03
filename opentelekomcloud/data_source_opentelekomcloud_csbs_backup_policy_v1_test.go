package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCSBSBackupPolicyV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1DataSourceID("data.opentelekomcloud_csbs_backup_policy_v1.csbs_policy"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_policy_v1.csbs_policy", "name", "csbs-policy"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_policy_v1.csbs_policy", "description", "test-code"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_policy_v1.csbs_policy", "provider_id", "fc4d5750-22e7-4798-8a46-f48f62c4c1da"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1DataSourceID(n string) resource.TestCheckFunc {
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

var testAccCSBSBackupPolicyV1DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "s2.medium.1"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
	name            = "csbs-policy"
  	description      = "test-code"
  	provider_id = "fc4d5750-22e7-4798-8a46-f48f62c4c1da"
  	common= {  }
  	resources = [{
    resource_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
    resource_type = "OS::Nova::Server"
    resource_name = "resource4"
  	}]
  	scheduled_operations = [{
    scheduled_period_name ="mybackup"
    enabled = true
    scheduled_period_description = "My backup policy"
    operation_type ="backup"
    max_backups = "20"
    pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  	}]
}
data "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {  
  id = "${opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1.id}"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
