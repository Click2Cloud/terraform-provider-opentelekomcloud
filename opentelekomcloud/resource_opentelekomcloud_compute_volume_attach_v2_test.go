package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

// PASS
func TestAccComputeV2VolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2VolumeAttach_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("opentelekomcloud_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

// PASS
func TestAccComputeV2VolumeAttach_device(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2VolumeAttach_device,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("opentelekomcloud_compute_volume_attach_v2.va_1", &va),
					//testAccCheckComputeV2VolumeAttachDevice(&va, "/dev/vdc"),
				),
			},
		},
	})
}

// PASS
func TestAccComputeV2VolumeAttach_timeout(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2VolumeAttach_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("opentelekomcloud_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckComputeV2VolumeAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_volume_attach_v2" {
			continue
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2VolumeAttachExists(n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceId || found.VolumeID != volumeId {
			return fmt.Errorf("VolumeAttach not found")
		}

		*va = *found

		return nil
	}
}

func testAccCheckComputeV2VolumeAttachDevice(
	va *volumeattach.VolumeAttachment, device string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if va.Device != device {
			return fmt.Errorf("Requested device of volume attachment (%s) does not match: %s",
				device, va.Device)
		}

		return nil
	}
}

var testAccComputeV2VolumeAttach_basic = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  volume_id = "${opentelekomcloud_blockstorage_volume_v2.volume_1.id}"
}
`, OS_NETWORK_ID)

var testAccComputeV2VolumeAttach_device = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  volume_id = "${opentelekomcloud_blockstorage_volume_v2.volume_1.id}"
  device = "/dev/vdc"
}
`, OS_NETWORK_ID)

var testAccComputeV2VolumeAttach_timeout = fmt.Sprintf(`
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
  volume_id = "${opentelekomcloud_blockstorage_volume_v2.volume_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_NETWORK_ID)
