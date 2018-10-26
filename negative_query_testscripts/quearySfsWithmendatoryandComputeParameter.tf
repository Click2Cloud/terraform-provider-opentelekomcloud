resource "opentelekomcloud_vpc_v1" "vpc_sfs011" {
  name = "vpc_sfs_test11"
  cidr = "192.168.0.0/16"
}
resource "opentelekomcloud_sfs_file_system_v2" "Share_file_011" {
  size = 40
  name = "test-share-11"
  access_to = "${opentelekomcloud_vpc_v1.vpc_sfs011.id}"
  access_level = "rw"
  region = "eu-de"
  share_proto = "NFS"
  description = "Ceate Sfs with optional parameters"
  is_public = false
  access_type = "cert"
}
data "opentelekomcloud_sfs_file_system_v2" "Share_file_012" {
  name = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.name}"
  id = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.id}"
  status = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.status}"
  availability_zone = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.availability_zone}"
  size = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.size}"
  host = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.host}"
  share_proto = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.share_proto}"
}
output "region" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.region}"
}
output "id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.id}"
}
output "availability_zone" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.availability_zone}"
}
output "size" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.size}"
}
/*output "share_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_type}"
}*/
output "status" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.status}"
}
output "description" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.description}"
}
output "host" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.host}"
}
output "is_public" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.is_public}"
}
output "name" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.name}"
}
output "share_proto" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.share_proto}"
}
/*output "volume_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.volume_type}"
}
output "export_location" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.export_location}"
}
output "metadata" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.metadata}"
}*/
/*output "export_locations" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.export_locations}"
}*/
output "access_level" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.access_level}"
}
/*output "state" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.state}"
}*/
output "share_access_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.share_access_id}"
}
output "access_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.access_type}"
}
output "access_to" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_011.access_to}"
}
/*output "mount_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.mount_id}"
}
output "share_instance_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.share_instance_id}"
}

output "preferred" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_004.preferred}"
}*/
