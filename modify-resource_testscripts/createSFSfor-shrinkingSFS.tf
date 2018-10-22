resource "opentelekomcloud_vpc_v1" "vpc_sfs013" {
  name = "vpc_sfs_test13"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_013" {
  size = 60
  name = "test-share-13"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs013.id}"
  access_level = "rw"
}