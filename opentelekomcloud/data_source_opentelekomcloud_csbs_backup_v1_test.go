package opentelekomcloud
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccCSBSBackupV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupV1DataSourceID("data.opentelekomcloud_csbs_backup_v1.csbs"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_v1.csbs", "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_v1.csbs", "description", "test-code"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_csbs_backup_v1.csbs", "resource_type", "OS::Nova::Server"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupV1DataSourceID(n string) resource.TestCheckFunc {
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

var testAccCSBSBackupV1DataSource_basic = `
resource "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test"
  description      = "test-code"
  resource_id = "92cc41e5-a761-4828-9b19-247076aa4e55"
  resource_type = "OS::Nova::Server"
}
data "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_id = "${opentelekomcloud_csbs_backup_v1.csbs.id}"
}
`
