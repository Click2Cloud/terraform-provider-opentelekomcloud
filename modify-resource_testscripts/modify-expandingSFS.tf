resource "${opentelekomcloud_sfs_file_system_v2.Share_file_012.id}" {
  size = 50
  name = "test-share11"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs012.id}"
  access_level = "rw"
}