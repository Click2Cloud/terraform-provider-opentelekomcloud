resource "opentelekomcloud_vpc_v1" "vpc_sfs011" {
  name = "vpc_sfs_test11"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_011" {
  size = 40
  name = "test-share-11"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs011.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_012" {
  name = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.name}"
  id = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.id}"
  availability_zone = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.availability_zone}"
  size = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.size}"
  host = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.host}"
  share_proto = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.share_proto}"

}