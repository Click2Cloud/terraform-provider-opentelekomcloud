resource "opentelekomcloud_vpc_v1" "vpc_sfs008" {
  name = "vpc_sfs_test8"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_008" {
  name = "test-share-08"
  availability_zone = "	eu-de_terraform"
}
