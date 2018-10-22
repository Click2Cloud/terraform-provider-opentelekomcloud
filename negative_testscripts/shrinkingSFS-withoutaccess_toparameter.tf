resource "opentelekomcloud_vpc_v1" "vpc_sfs014" {
  name = "vpc_sfs_test14"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_014" {
  size = 50
  name = "test-share-14"
  access_level = "rw"
}