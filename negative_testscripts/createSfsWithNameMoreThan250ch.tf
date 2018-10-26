resource "opentelekomcloud_vpc_v1" "vpc_sfs007" {
  name = "vpc_sfs_test7"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_007" {
  size = 100
  name = "test-share-111nvdsndskdscd45ddvieiffk4grlaevvsajdjjefjjofrrr5g7rg6r6sfdsa54f6fasf4e4fd44ggggg54sv4s5test-share-111nvdsndskdscd45ddvieiffk4grlaevvsajdjjefjjofrrr5g7rg6r6sfdsa54f6fasf4e4fd44ggggg54sv4s5ffk4grlaevvsajdjjefjjofrrr5g7rg6r6sfdsa54f6fasf4e4fd44ggggg54sv4s5"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs007.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}
