resource "opentelekomcloud_vpc_v1" "vpc_sfs002" {
  name = "vpc_sf_test2"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_002" {
  size = 100
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs002.id}"
  access_level = "rw"
}