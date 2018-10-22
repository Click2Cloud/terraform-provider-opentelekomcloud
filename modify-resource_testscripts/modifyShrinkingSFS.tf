resource "opentelekomcloud_sfs_file_system_v2.Share_file_013" {
  size = 40
  name = "test-share11"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs013.id}"
  access_level = "rw"
}