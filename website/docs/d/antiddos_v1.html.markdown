---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_antiddos_v1"
sidebar_current: "docs-opentelekomcloud-datasource-antiddos-v1"
description: |-
  Provides status of a specific EIP.
---

# Data Source: opentelekomcloud_antiddos_v1

The OpenTelekomCloud Antiddos data source allows to query the status of EIP, regardless whether an EIP has been bound to an Elastic Cloud Server (ECS) or not.

## Example Usage


```hcl
variable "id" { }

data "opentelekomcloud_antiddos_v1" "antiddos" {
  floating_ip_id = "${opentelekomcloud_antiddos_v1.antiddos_1.id}"
}

```

## Argument Reference
The following arguments are supported:

* `floating_ip_id` - (Optional) Specifies the id of an eip.

* `floating_ip_address` - (Optional) Specifies the floating ip address.

* `status` - (Optional) Specifies the defense status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_type` - Specifies the EIP type.

* `period_start` - Provides the Start time.

* `bps_attack` - Specifies the Attack traffic in (bit/s).

* `bps_in` - Specifies the inbound traffic in (bit/s).

* `total_bps` - Specifies the total traffic.

* `pps_in` - Specifies the inbound packet rate (number of packets per second).

* `pps_attack` - Specifies the attack packet rate (number of packets per second).

* `total_pps` - Specifies the total packet rate.

* `start_time` - Specifies the start time of cleaning and blackhole event.

* `end_time` - Specifies the end time of cleaning and blackhole event.

* `traffic_cleaning_status` - Specifies the traffic cleaning status.

* `trigger_bps` - Specifies the traffic at the triggering point.

* `trigger_pps` - Specifies the packet rate at the triggering point.

* `trigger_http_pps` - Specifies the HTTP request rate at the triggering point.

