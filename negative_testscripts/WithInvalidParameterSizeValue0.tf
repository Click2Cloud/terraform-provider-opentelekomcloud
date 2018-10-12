resource "opentelekomcloud_sfs_file_system_v2" "Share_file_001" {
  size = 0
  name = "test-share-01"
  access_to = "34af42db-654f-4128-9126-9ef312bdd2d4	",
  access_type = "cert",
  access_level = "rw"
}
output "availability_zone" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.availability_zone}"
}
output "size" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.size}"
}
output "share_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.share_type}"
}
output "project_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.project_id}"
}
output "status" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.status}"
}
output "description" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.description}"
}
output "host" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.host}"
}
output "is_public" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.is_public}"
}
output "name" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.name}"
}
output "share_proto" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.share_proto}"
}
output "volume_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.volume_type}"
}
output "export_location" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.export_location}"
}
output "metadata" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.metadata}"
}
output "export_locations" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.export_locations}"
}
output "access_level" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.access_level}"
}
output "state" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.state}"
}
output "share_access_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.share_access_id}"
}
output "access_type" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.access_type}"
}
output "access_to" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.access_to}"
}
output "mount_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.mount_id}"
}
output "share_instance_id" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.share_instance_id}"
}
output "preferred" {
  value = "${opentelekomcloud_sfs_file_system_v2.Share_file_001.preferred}"
}
