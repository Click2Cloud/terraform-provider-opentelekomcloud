resource "opentelekomcloud_sfs_file_system_v2.Share_file_011" {
  size = 60
  name = "test-share1111"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs011.id}"
  access_level = "rw"
}
