package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/policies"
	"testing"
)

func TestAccVBSBackupPolicyV2_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "backup_policy_name", "policy_001"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "backup_policy_name", "policy_002"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
		},
	})
}

func TestAccVBSBackupPolicyV2_timeout(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("opentelekomcloud_vbs_backup_policy_v2.vbs", &policy),
				),
			},
		},
	})
}

func testAccVBSBackupPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vbs_backup_policy_v2" {
			continue
		}

		_, err := policies.List(vbsClient, policies.ListOpts{PolicyID: rs.Primary.ID})
		if err == nil {
			return fmt.Errorf("Share File still exists")
		}
	}

	return nil
}

func testAccVBSBackupPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating opentelekomcloud vbs client: %s", err)
		}

		policyList, err := policies.List(vbsClient, policies.ListOpts{PolicyID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.PolicyID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = found

		return nil
	}
}

var testAccVBSBackupPolicyV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
  backup_policy_name = "policy_001"
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
`)

var testAccVBSBackupPolicyV2_update = fmt.Sprintf(`
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
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
`)

var testAccVBSBackupPolicyV2_timeout = fmt.Sprintf(`
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
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

  timeouts {
    create = "5m"
    delete = "5m"
  }
}`)
