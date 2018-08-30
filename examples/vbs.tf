resource "opentelekomcloud_vbs_backup_policy_v2" "vbs" {
  name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }]
}

data "opentelekomcloud_vbs_backup_policy_v2" "policies" {
  id = "${opentelekomcloud_vbs_backup_policy_v2.vbs.id}"
}
