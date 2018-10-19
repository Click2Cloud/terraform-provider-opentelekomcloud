resource "opentelekomcloud_vpc_v1" "vpc_sfs009" {
  name = "vpc_sfs_test8"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_009" {
  size = 60
  name = "test-share-101"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs009.id}"
  access_level = "rw"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_004" {
  availability_zone = "${opentelekomcloud_sfs_file_system_v2.Share_file_009.availability_zone}"
  size = "${opentelekomcloud_sfs_file_system_v2.Share_file_009.size}"
  host = "${opentelekomcloud_sfs_file_system_v2.Share_file_009.host}"
  share_proto = "${opentelekomcloud_sfs_file_system_v2.Share_file_009.share_proto}"
}