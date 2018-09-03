package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/backups"
)

func TestAccOTCVBSBackupV2_basic(t *testing.T) {
	var config backups.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("opentelekomcloud_vbs_backup_v2.backup_1", &config),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "name", "opentelekomcloud-backup"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "description", "Backup_Demo"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_v2.backup_1", "status", "available"),
				),
			},
		},
	})
}

func TestAccOTCVBSBackupV2_timeout(t *testing.T) {
	var config backups.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("opentelekomcloud_vbs_backup_v2.backup_1", &config),
				),
			},
		},
	})
}

func testAccCheckVBSBackupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vbs_backup_v2" {
			continue
		}

		_, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VBS backup still exists")
		}
	}

	return nil
}

func testAccCheckVBSBackupV2Exists(n string, configs *backups.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud vbs client: %s", err)
		}

		found, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("VBS backup not found")
		}

		*configs = *found

		return nil
	}
}

const testAccVBSBackupV2_basic = `
resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id = "b02b11ea-4eab-4bcb-96b7-9c872adfdafc"
  name = "opentelekomcloud-backup"
  description = "Backup_Demo"
}
`

const testAccVBSBackupV2_timeout = `
resource "opentelekomcloud_vbs_backup_v2" "backup_1" {
  volume_id = "b02b11ea-4eab-4bcb-96b7-9c872adfdafc"
  name = "opentelekomcloud-backup"
  description = "Backup_Demo"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`