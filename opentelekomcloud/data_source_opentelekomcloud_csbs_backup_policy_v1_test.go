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

var testAccCSBSBackupPolicyV1DataSource_basic = `
resource "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {
  name            = "csbs-policy"
  description      = "test-code"
  provider_id = "fc4d5750-22e7-4798-8a46-f48f62c4c1da"
  common= {  }
  resources = [{
    id = "9422f270-6fcf-4ba2-9319-a007f2f63a8e"
    type = "OS::Nova::Server"
    r_name = "resource4"
  }]
  scheduled_operations = [{
    so_name ="mybackupp"
    enabled = true
    so_description = "My backup policyy"
    operation_type ="backup"
    max_backups = "3"
    pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }]
}
data "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {  
  id = "${opentelekomcloud_csbs_backup_policy_v1.csbs_policy.id}"
}
`
