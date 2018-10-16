resource "opentelekomcloud_vpc_v1" "vpc_sfs007" {
  name = "vpc_sfs_test6"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_007" {
  name = "test-share-07"
  availability_zone = "	eu-de_terraform"
}
