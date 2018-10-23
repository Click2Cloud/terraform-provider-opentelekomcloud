resource "opentelekomcloud_vpc_v1" "vpc_sfs017" {
  name = "vpc_sfs_test17"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_017" {
  size = 30
  name = "test-share-17"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs017.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}