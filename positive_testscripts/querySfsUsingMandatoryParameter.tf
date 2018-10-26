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
  access_type = "cert"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_005" {
  name = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.name}"
  id = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.id}"
  status = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.status}"
}
output "region" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.region}"
}
output "id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.id}"
}
output "availability_zone" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.availability_zone}"
}
output "size" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.size}"
}
output "share_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_type}"
}
output "status" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.status}"
}
output "description" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.description}"
}
output "host" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.host}"
}
output "is_public" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.is_public}"
}
output "name" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.name}"
}
output "share_proto" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_proto}"
}
output "volume_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.volume_type}"
}
output "export_location" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.export_location}"
}
output "metadata" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.metadata}"
}
output "export_locations" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.export_locations}"
}
output "access_level" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.access_level}"
}
output "state" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.state}"
}
output "share_access_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_access_id}"
}
output "access_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.access_type}"
}
output "access_to" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.access_to}"
}
output "mount_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.mount_id}"
}
output "share_instance_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_instance_id}"
}
output "preferred" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.preferred}"
}