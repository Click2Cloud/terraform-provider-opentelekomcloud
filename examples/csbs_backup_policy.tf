resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name                = "${var.project}-policy"
  description         = "test-code"
  provider_id         = "fc4d5750-22e7-4798-8a46-f48f62c4c1da"
  common              = {  }
  resources = [{
    resource_id       = "${opentelekomcloud_compute_instance_v2.webserver.id}"
    resource_type     = "OS::Nova::Server"
    resource_name     = "resource1"
  }]
  scheduled_operations = [{
    scheduled_period_name         = "mybackup"
    enabled                       =  true
    scheduled_period_description  =  "My backup policy"
    operation_type                =  "backup"
    max_backups                   =  15
    permanent                     =  true
    pattern                       =  "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }]
}
