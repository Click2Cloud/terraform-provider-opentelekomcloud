resource "opentelekomcloud_vpc_v1" "vpc_sfs004" {
  name = "vpc_sfs_test4"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_004" {
  size = 40
  name = "test-share-04"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs004.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  availability_zone = "eu-de_terraform"
  access_type = "cert"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_005" {
  name = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.name}"
  id = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.id}"

}