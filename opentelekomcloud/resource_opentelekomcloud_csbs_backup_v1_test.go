package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"log"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
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
						"opentelekomcloud_csbs_backup_v1.csbs", "backup_name", "csbs-test1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_v1.csbs", "description", "test-code"),
				),
			},
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
	backupClient, err := config.backupV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud backup client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_v1" {
			continue
		}

		_, err := backup.List(backupClient, backup.ListOpts{CheckpointId: rs.Primary.Attributes["backup_record_id"]})
		if err != nil {
			return fmt.Errorf("Backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Backup not found: %s", n)
		}

		if rs.Primary.Attributes["backup_record_id"] == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		backupClient, err := config.backupV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating opentelekomcloud csbs client: %s", err)
		}

		backupList, err := backup.List(backupClient, backup.ListOpts{CheckpointId: rs.Primary.Attributes["backup_record_id"]})
		if err != nil {
			return err
		}
		found := backupList[0]
		log.Printf("[DEBUG] found : %s", found)
		if found.CheckpointId != rs.Primary.Attributes["backup_record_id"] {
			return fmt.Errorf("backup  not found")
		}

		*backups = found

		return nil
	}
}

var testAccCSBSBackupV1_basic = fmt.Sprintf(`
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
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccCSBSBackupV1_timeout = fmt.Sprintf(`
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
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
