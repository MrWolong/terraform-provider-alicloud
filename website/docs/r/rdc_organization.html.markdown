---
subcategory: "Apsara Devops(RDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_rdc_organization"
sidebar_current: "docs-alicloud-resource-rdc-organization"
description: |-
  Provides a Alicloud RDC Organization resource.
---

# alicloud_rdc_organization

Provides a RDC Organization resource.

For information about RDC Organization and how to use it, see [What is Organization](https://www.alibabacloud.com/help/en/yunxiao/product-overview/what-is-cloud-effect).

-> **NOTE:** Available since v1.137.0.

-> **DEPRECATED:** This resource has been deprecated from version `1.238.0`.

## Example Usage

Basic Usage

```terraform
resource "alicloud_rdc_organization" "example" {
  organization_name = "example_value"
  source            = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `desired_member_count` - (Optional) The desired member count.
* `organization_name` - (Required, ForceNew, ForceNew) Company name.
* `real_pk` - (Optional) User pk, not required, only required when the ak used by the calling interface is inconsistent with the user pk
* `source` - (Required) This is organization source information

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Organization.

## Import

RDC Organization can be imported using the id, e.g.

```shell
$ terraform import alicloud_rdc_organization.example <id>
```
