package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "description", "test-code"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "provider_id", "fc4d5750-22e7-4798-8a46-f48f62c4c1da"),
				),
			},
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy-update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "description", "test-code-update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", "provider_id", "fc4d5750-22e7-4798-8a46-f48f62c4c1da"),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	policyClient, err := config.backupV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud backup policy client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_csbs_backup_policy_v1" {
			continue
		}

		_, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("backup policy still exists")
		}
	}

	return nil
}

func testAccCheckCSBSBackupPolicyV1Exists(n string, policy *policies.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		policyClient, err := config.backupV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating opentelekomcloud backup policy client: %s", err)
		}

		found, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = *found

		return nil
	}
}

var testAccCSBSBackupPolicyV1_basic = fmt.Sprintf(`
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
	name            = "backup-policy"
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
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_update = fmt.Sprintf(`
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
	name            = "backup-policy-update"
  	description      = "test-code-update"
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
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_timeout = fmt.Sprintf(`
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
	name            = "backup-policy"
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
	timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
