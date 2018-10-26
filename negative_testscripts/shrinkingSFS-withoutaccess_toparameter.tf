resource "opentelekomcloud_vpc_v1" "vpc_sfs015" {
  name = "vpc_sfs_test15"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_015" {
  size = 50
  name = "test-share-14"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}