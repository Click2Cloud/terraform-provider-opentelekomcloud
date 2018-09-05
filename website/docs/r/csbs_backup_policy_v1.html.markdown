---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: resource_opentelekomcloud_csbs_backup_policy_v1"
sidebar_current: "docs-opentelekomcloud-resource-csbs-backup-policy-v1"
description: |-
  Provides an OpenTelekomCloud Backup Policy of Resources.
---

# opentelekomcloud_csbs_backup_policy_v1

Provides an OpenTelekomCloud Backup Policy of Resources.

## Example Usage

 ```hcl
 variable "name" { }
 variable "description" { }
 variable "resource_id" { }
 variable "resource_name" { }
 variable "scheduled_period_name" { }
 
 resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
   name            = "${var.name}"
   description      = "${var.description}"
   provider_id = "fc4d5750-22e7-4798-8a46-f48f62c4c1da"
   common= {  }
   resources = [{
     resource_id = "${var.resource_id}"
     resource_type = "OS::Nova::Server"
     resource_name = "${var.resource_name}"
   }]
   scheduled_operations = [{
     scheduled_period_name ="${var.scheduled_period_name}"
     enabled = true
     operation_type ="backup"
     pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
   }]
 }

 ```
## Argument Reference
The following arguments are supported:

* `name` - (Required) Specifies the name of backup policy. The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional) Backup policy description. The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

* `provider_id` - (Required) Specifies backup provider ID.

* `common` - (Optional) General backup policy parameters, which are blank by default.

**scheduled_operations** **- (Optional)** Specifies Scheduling period. A backup policy has only one backup period.

* `scheduled_period_name` - (Optional) Specifies Scheduling period name.The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).
    
* `scheduled_period_description` - (Optional) Specifies Scheduling period description.The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

* `enabled` - (Optional) Specifies whether the scheduling period is enabled.

* `max_backups` - (Optional) Specifies maximum number of backups that can be automatically created for a backup object.

* `retention_duration_days` - (Optional) Specifies duration of retaining a backup, in days.

* `permanent` - (Optional) Specifies whether backups are permanently retained.

* `plan_id` - (Optional) Specifies backup policy ID.

* `pattern` - (Required) Specifies Scheduling policy of the scheduler.

* `operation_type` - (Required) Specifies Operation type, which can be backup.
**resources** **- (Optional)** Backup object List. The list can be blank.

* `resource_id` - (Required) Specifies the ID of the object to be backed up.
    
* `resource_type` - (Required) Entity object type of the backup object. If the type is VMs, the value is **OS::Nova::Server**.

* `resource_name` - (Required) Specifies backup object name.
**tags** **- (Optional)** List of tags. Keys in this list must be unique.

* `key` - (Required) Tag key. It cannot be an empty string.
    
* `value` - (Required) Tag value. It can be an empty string.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `status` - Status of Backup Policy.

* `id` - Backup Policy ID.

**scheduled_operations** - Backup plan information

* `scheduled_period_id` -  Specifies Scheduling period ID.

* `scheduler_id` -  Specifies Scheduler ID.

* `scheduler_name` -  Specifies Scheduler name.

* `scheduler_type` -  Specifies Scheduler type.


## Import

Backup Policy can be imported using  `id`, e.g.

```
$ terraform import opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```




