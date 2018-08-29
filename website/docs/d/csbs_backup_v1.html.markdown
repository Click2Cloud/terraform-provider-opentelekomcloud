---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_csbs_backup_v1"
sidebar_current: "docs-opentelekomcloud-datasource-csbs-backup-v1"
description: |-
  Provides details about a specific Backup.
---

# Data Source: opentelekomcloud_csbs_backup_v1

The OpenTelekomCloud CSBS Backup data source allows access of backup resources.

## Example Usage


```hcl
variable "backup_name" { }

data "opentelekomcloud_csbs_backup_v1" "csbs" {
  backup_name = "${var.backup_name}" 
}
```

## Argument Reference
The following arguments are supported:

* `backup_id` - (Optional) Specifies the ID of backup.

* `backup_name` - (Optional) Specifies the backup name.

* `status` - (Optional) Specifies the backup status.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `backup_record_id` - Specifies the backup record ID.

* `resource_id` - Specifies the backup object ID.

* `description` - Provides the backup description.

* `resource_type` - Specifies the type of backup objects.

* `auto_trigger` - Specifies whether automatic trigger is enabled.

* `average_speed` - Specifies average speed.

* `copy_from` - This parameter is left blank by default.

* `copy_status` - The default value is **na**.

* `fail_op` - Specifies the type of the failed operation i.e. **backup, restore, delete**.

* `fail_reason` - Specifies the description of the failure cause,

* `image_type` - Specifies the backup type. Default value is **backup**

* `incremental` - Specifies whether incremental backup is used.

* `progress` - Specifies the progress.

* `resource_az` - Specifies the AZ to which the backup resource belongs.

* `resource_name` - Specifies the backup object name.

* `size` - Specifies the backup capacity.

* `space_saving_ratio` - Specifies the space saving rate.

* `supported_restore_mode` - Specifies the restoration mode.

* `support_lld` - Specifies whether to allow lazyloading for fast restoration.

* `cloudservicetype` - Specifies the ECS type.

* `disk` - Specifies the system disk size corresponding to the ECS specifications.

* `imagetype` - Specifies the image type. The value can be: gold: public image, private: private image,  market: market image.

* `ram` - Specifies the memory size of the ECS, in MB.

* `vcpus` - Specifies the cpu cores corresponding to the ECS.

* `eip` - Specifies the Elastic IP address of the ECS.

* `private_ip` - Specifies the internal IP address of the ECS.
