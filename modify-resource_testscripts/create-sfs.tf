resource "opentelekomcloud_vpc_v1" "vpc_sfs011" {
  name = "vpc_sfs_test11"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_011" {
  size = 60
  name = "test-share-11"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs011.id}"
  access_level = "rw"
}