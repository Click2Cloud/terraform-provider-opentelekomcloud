---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cts_tracker_v1"
sidebar_current: "docs-opentelekomcloud-datasource-cts-tracker-v1"
description: |-
  Allows you to collect, store, and query cloud resource operation records and use these records for security analysis, compliance auditing, resource tracking, and fault locating.
---

# Data Source: opentelekomcloud_cts_tracker_v1

The OpenTelekomCloud CTS Tracker data source allows access of Cloud Tracker.

## Example Usage


```hcl
variable "bucket_name" { }

data "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name = "${var.bucket_name}"
}

```

## Argument Reference
The following arguments are supported:

* `tracker_name` - (Optional) Specifies the tracker name. 

* `bucket_name` - (Optional) Specifies the OBS bucket name.

* `file_prefix_name` - (Optional) Specifies the prefix of a log that needs to be stored in an OBS bucket. 

* `status` - (Optional) Specifies the status of a tracker. 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `smn` block supports the following arguments:

    * `is_support_smn` - Specifies whether SMN is supported.
    
    * `topic_id` - 	Specifies the theme of the SMN service.

    * `operations` - Specifies trigger conditions for sending a notification

    * `is_send_all_key_operation` - When the value is false, operations cannot be left empty.

    * `need_notify_user_list` - Specifies the users using the login function. When these users log in, notifications will be sent.

    