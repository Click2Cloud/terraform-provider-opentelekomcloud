resource "opentelekomcloud_sfs_file_system_v2.Share_file_019" {
  size = 50
  name = "test-share19"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs017.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}