resource "opentelekomcloud_vpc_v1" "vpc_sfs013" {
  name = "vpc_sfs_test13"
  cidr="192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_013" {
  size = 60
  name = "test-share-13"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs013.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_012" {
  availability_zone = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.availability_zone}"
  size = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.size}"
  host = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.host}"
  share_proto = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_proto}"
  region = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.region}"
  project_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.project_id}"
  status = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.status}"
  description = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.description}"
  is_public = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.is_public}"
  volume_type = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.volume_type}"
  export_location = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.export_location}"
  metadata = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.metadata}"
  access_level = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_level}"
  state = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.state}"
  share_access_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_access_id}"
  access_type = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_type}"
  access_to = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_to}"
  mount_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.mount_id}"
  share_instance_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_instance_id}"
  preferred = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.preferred}"
}
output "region" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.region}"
}
output "id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.id}"
}
output "availability_zone" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.availability_zone}"
}
output "size" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.size}"
}
output "share_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_type}"
}
output "status" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.status}"
}
output "description" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.description}"
}
output "host" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.host}"
}
output "is_public" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.is_public}"
}
output "name" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.name}"
}
output "share_proto" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_proto}"
}
output "volume_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.volume_type}"
}
output "export_location" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.export_location}"
}
output "metadata" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.metadata}"
}
output "export_locations" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.export_locations}"
}
output "access_level" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_level}"
}
output "state" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.state}"
}
output "share_access_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_access_id}"
}
output "access_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_type}"
}
output "access_to" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.access_to}"
}
output "mount_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.mount_id}"
}
output "share_instance_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.share_instance_id}"
}
output "preferred" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_013.preferred}"
}