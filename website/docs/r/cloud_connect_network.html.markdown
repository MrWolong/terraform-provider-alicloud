---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network"
sidebar_current: "docs-alicloud-resource-cloud-connect-network"
description: |-
  Provides a Alicloud Cloud Connect Network resource.
---

# alicloud_cloud_connect_network

Provides a cloud connect network resource. Cloud Connect Network (CCN) is another important component of Smart Access Gateway. It is a device access matrix composed of Alibaba Cloud distributed access gateways. You can add multiple Smart Access Gateway (SAG) devices to a CCN instance and then attach the CCN instance to a Cloud Enterprise Network (CEN) instance to connect the local branches to the Alibaba Cloud.

For information about cloud connect network and how to use it, see [What is Cloud Connect Network](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createcloudconnectnetwork).

-> **NOTE:** Available since v1.59.0.

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_connect_network&exampleId=0d167127-4de5-23ad-4668-fd36ef279ab0a70ff28c&activeTab=example&spm=docs.r.cloud_connect_network.0.0d1671274d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
resource "alicloud_cloud_connect_network" "default" {
  name        = var.name
  description = var.name
  cidr_block  = "192.168.0.0/24"
  is_default  = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CCN instance. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.
* `description` - (Optional) The description of the CCN instance. The description can contain 2 to 256 characters. The description must start with English letters, but cannot start with http:// or https://.
* `cidr_block` - (Optional) The CidrBlock of the CCN instance. Defaults to null.
* `is_default` - (Required) Created by default. If the client does not have ccn in the binding, it will create a ccn for the user to replace.


## Attributes Reference

The following attributes are exported:

* `id` - The CcnId of the CCN instance. For example "ccn-xxx".

## Import

The cloud connect network instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_connect_network.example ccn-abc123456
```

