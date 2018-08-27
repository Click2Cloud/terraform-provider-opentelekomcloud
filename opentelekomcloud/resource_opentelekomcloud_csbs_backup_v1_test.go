package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
	"testing"
	"log"
)

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("opentelekomcloud_csbs_backup_v1.csbs", &backups),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "description", "test-code"),
				),
			},
			/*resource.TestStep{
				Config: testAccCSBSBackupV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("opentelekomcloud_csbs_backup_v1.csbs", &backups),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "backup_policy_name", "policy_002"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "status", "ON"),
				),
			},*/
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("opentelekomcloud_csbs_backup_v1.csbs", &backups),
				),
			},
		},
	})
}

func testAccCSBSBackupV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.backupV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_v1" {
			continue
		}

		_, err := backup.List(vbsClient, backup.ListOpts{CheckpointId: rs.Primary.ID})
		if err == nil {
			return fmt.Errorf("Backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vbsClient, err := config.backupV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating opentelekomcloud csbs client: %s", err)
		}

		policyList, err := backup.List(vbsClient, backup.ListOpts{CheckpointId: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.CheckpointId != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		//*backups = found
		log.Printf("found : %#v", found)

		return nil
	}
}

var testAccCSBSBackupV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test"
  description      = "test-code"
  resource_id = "92cc41e5-a761-4828-9b19-247076aa4e55"
  resource_type = "OS::Nova::Server"
}
`)

/*var testAccCSBSBackupV1_update = fmt.Sprintf(`
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_policy_name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "Y"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }] 
}
`)*/

var testAccCSBSBackupV1_timeout = fmt.Sprintf(`
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test"
  description      = "test-code"
  resource_id = "92cc41e5-a761-4828-9b19-247076aa4e55"
  resource_type = "OS::Nova::Server"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}`)
