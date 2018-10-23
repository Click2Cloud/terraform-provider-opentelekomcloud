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
data "opentelekomcloud_sfs_file_system_v2" "Share_file_014" {
  availability_zone = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.availability_zone}"
  size = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.size}"
  host = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.host}"
  share_proto = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.share_proto}"
  region = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.region}"
  project_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.project_id}"
  status = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.status}"
  description = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.description}"
  is_public = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.is_public}"
  volume_type = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.volume_type}"
  export_location = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.export_location}"
  metadata = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.metadata}"
  access_level = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.access_level}"
  state = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.state}"
  share_access_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.share_access_id}"
  access_type = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.access_type}"
  access_to = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.access_to}"
  mount_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.mount_id}"
  share_instance_id = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.share_instance_id}"
  preferred = "${opentelekomcloud_sfs_file_system_v2.Share_file_014.preferred}"
}