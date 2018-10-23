resource "opentelekomcloud_vpc_v1" "vpc_sfs03" {
  name = "vpc_sfs_test03"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_03" {
  size = 60
  name = "test-share-03"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs03.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}
