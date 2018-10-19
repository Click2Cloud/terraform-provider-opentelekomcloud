resource "opentelekomcloud_vpc_v1" "vpc_sfs010" {
  name = "vpc_sfs_test10"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_010" {
  size = 60
  name = "test-share-10"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs010.id}"
  access_level = "rw"
}
