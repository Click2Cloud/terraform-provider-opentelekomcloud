resource "opentelekomcloud_vpc_v1" "vpc_sfs012" {
  name = "vpc_sfs_test12"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_012" {
  size = 30
  name = "test-share-12"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs012.id}"
  access_level = "rw"
}