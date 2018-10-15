resource "opentelekomcloud_vpc_v1" "vpc_sfs006" {
  name = "vpc_sfs_test5"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_006" {
  size = 0
  name = "test-share-06"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs001.id}"
  access_type = "cert",
  access_level = "rw"
}
