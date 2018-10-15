resource "opentelekomcloud_vpc_v1" "vpc_sfs003" {
  name = "vpc_sfs_test3"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_003" {
  size = 40
  name = "test-share-101"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs003.id}"
  access_level = "rw"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_004" {
  name = "${opentelekomcloud_sfs_file_system_v2.Share_file_003.name}"
  id = "${opentelekomcloud_sfs_file_system_v2.Share_file_003.id}"

}