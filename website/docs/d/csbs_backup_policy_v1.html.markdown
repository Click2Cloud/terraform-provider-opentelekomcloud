---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_csbs_backup_policy_v1"
sidebar_current: "docs-opentelekomcloud-datasource-csbs-backup-policy-v1"
description: |-
  Provides details about a specific Backup Policy.
---

# Data Source: opentelekomcloud_csbs_backup_policy_v1

The OpenTelekomCloud CSBS Backup Policy data source allows access of backup Policy resources.

## Example Usage


```hcl
variable "id" { }

data "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {
  id = "${var.id}" 
}

```

## Argument Reference
The following arguments are supported:

* `id` - (Optional) Specifies the ID of backup policy.

* `name` - (Optional) Specifies the backup policy name.

* `status` - (Optional) Specifies the backup policy status.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `description` - Specifies the backup policy description.

* `project_id` - Specifies the Tenant ID.

* `provider_id` - Provides the Backup provider ID.

* `parameters` - Specifies the parameters of a backup policy.

**scheduled_operations** - Specifies Scheduling period. A backup policy has only one backup period.

* `scheduled_period_name` - Specifies Scheduling period name.
    
* `scheduled_period_description` - Specifies Scheduling period description.

* `enabled` - Specifies whether the scheduling period is enabled.

* `max_backups` - Specifies maximum number of backups that can be automatically created for a backup object.

* `retention_duration_days` - Specifies duration of retaining a backup, in days.

* `permanent` - Specifies whether backups are permanently retained.

* `plan_id` - Specifies backup policy ID.

* `pattern` - Specifies Scheduling policy of the scheduler.

* `operation_type` - Specifies Operation type, which can be backup.

* `scheduled_period_id` -  Specifies Scheduling period ID.

* `scheduler_id` -  Specifies Scheduler ID.

* `scheduler_name` -  Specifies Scheduler name.

* `scheduler_type` -  Specifies Scheduler type.

**resources** - Backup object List. The list can be blank.

* `resource_id` - Specifies the ID of the object to be backed up.
    
* `resource_type` - Entity object type of the backup object. 

* `resource_name` - Specifies backup object name.

**tags** - List of tags. Keys in this list must be unique.

* `key` - Tag key. It cannot be an empty string.
    
* `value` - Tag value. It can be an empty string.
