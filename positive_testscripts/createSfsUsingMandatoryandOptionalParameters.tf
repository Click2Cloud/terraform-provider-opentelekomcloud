resource "opentelekomcloud_vpc_v1" "vpc_sfs001" {
  name = "vpc_sf_test1"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_001" {
  size = 100
  name = "test-share-01"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs001.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"

}