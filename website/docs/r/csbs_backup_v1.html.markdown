---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: resource_opentelekomcloud_csbs_backup_v1"
sidebar_current: "docs-opentelekomcloud-resource-csbs-backup-v1"
description: |-
  Provides an OpenTelekomCloud Backup of Resources.
---

# opentelekomcloud_csbs_backup_v1

Provides an OpenTelekomCloud Backup of Resources.

## Example Usage

 ```hcl
 variable "backup_name" { }
 variable "description" { }
 variable "resource_id" { }
 
 resource "opentelekomcloud_csbs_backup_v1" "backup_v1" {
   backup_name = "${var.backup_name}"
   description = "${var.description}"
   resource_id = "${var.resource_id}"
   resource_type = "OS::Nova::Server"
 }

 ```
## Argument Reference
The following arguments are supported:

* `backup_name` - (Optional) Name for the backup. The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-). Changing backup_name creates a new backup.

* `description` - (Optional) Backup description. The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<). Changing description creates a new backup.

* `resource_id` - (Required) ID of the target to which the backup is restored. Changing this creates a new backup.

* `resource_type` - (Required) Type of the target to which the backup is restored. The default value is **OS::Nova::Server** for an ECS. Changing this creates a new backup.

**tags** **- (Optional)** List of tags. Keys in this list must be unique. Changing this creates a new backup.

* `key` - (Required) Tag key. It cannot be an empty string.
    
* `value` - (Required) Tag value. It can be an empty string.
## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `status` - Status of Backup.

* `backup_record_id` - Backup record ID.

* `resource_graph` - Resource graph.

* `project_id` - Tenant Id of the Resource.

**protection_plan** - Backup plan information

* `id` -  Backup plan ID
    
* `name` -  Backup plan name
    
**resources** - List of Backup object.

* `id` - ID of the object to be backed up

* `type` - Backup object type. If the type is VMs, the value is **OS::Nova::Server**.

* `name` - Name of the backup object.


## Import

Backup can be imported using  `id`, e.g.

```
$ terraform import opentelekomcloud_csbs_backup_v1.backup_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```



