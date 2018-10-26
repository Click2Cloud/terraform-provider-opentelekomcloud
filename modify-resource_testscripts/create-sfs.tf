resource "opentelekomcloud_vpc_v1" "vpc_sfs016" {
  name = "vpc_sfs_test16"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_016" {
  size = 60
  name = "test-share-16"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs016.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}