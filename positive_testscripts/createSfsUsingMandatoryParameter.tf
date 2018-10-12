resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "terraform_provider_test"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_001" {
  size = 100
  name = "test-share-01"
  access_to = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  access_level = "rw"
}