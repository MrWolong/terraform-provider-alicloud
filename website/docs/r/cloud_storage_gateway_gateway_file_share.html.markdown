---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_file_share"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-file-share"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway File Share resource.
---

# alicloud_cloud_storage_gateway_gateway_file_share

Provides a Cloud Storage Gateway Gateway File Share resource.

For information about Cloud Storage Gateway Gateway File Share and how to use it, see [What is Gateway File Share](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewayfileshare).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway_file_share&exampleId=5328a1d8-77d9-a55f-1426-e628017cd1fe41938c73&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway_file_share.0.5328a1d877&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_oss_bucket" "default" {
  bucket = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  gateway_name             = var.name
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = var.name
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = true
  polling_interval        = 4500
  fe_limit                = 0
  backend_limit           = 0
  cache_mode              = "Cache"
  squash                  = "none"
  lag_period              = 5
}
```

## Argument Reference

The following arguments are supported:

* `access_based_enumeration` - (Optional) Whether to enable Windows ABE, the prime minister, need windowsAcl parameter is set to true in the entry into force of. Default value: `false`. **NOTE:** The attribute is valid when the attribute `protocol` is `SMB`. Gateway version >= 1.0.45 above support. 
* `backend_limit` - (Optional) The Max upload speed of the gateway file share. Unit: `MB/s`, 0 means unlimited. Value range: `0` ~ `1280`. Default value: `0`. **NOTE:** at the same time if you have to limit the maximum write speed, maximum upload speed is no less than the maximum write speed. 
* `browsable` - (Optional) The whether browsable of the gateway file share (that is, in the network neighborhood of whether you can find). The attribute is valid when the attribute `protocol` is `SMB`. Default value: `true`.
* `cache_mode` - (Optional, ForceNew) The set up gateway file share cache mode. Valid values: `Cache` or `Sync`. `Cache`: cached mode. `Sync`: replication mode are available. Default value: `Cache`.
* `direct_io` - (Optional, ForceNew) File sharing Whether to enable DirectIO (direct I/O mode for data transmission). Default value: `false`.
* `download_limit` - (Optional) The maximum download speed of the gateway file share. Unit: `MB/s`. `0` means unlimited. Value range: `0` ~ `1280`. **NOTE:** only in copy mode and enable download file data can be set. only when the shared opens the reverse synchronization or acceded to by the speed synchronization Group when, this parameter will not take effect. Gateway version >= 1.3.0 above support. 
* `fast_reclaim` - (Optional) The whether to enable Upload optimization of the gateway file share, which is suitable for data pure backup migration scenarios. Default value: `false`. **NOTE:** Gateway version >= 1.0.39 above support. 
* `fe_limit` - (Optional) The maximum write speed of the gateway file share. Unit: `MB/s`, `0` means unlimited. Value range: `0` ~ `1280`. Default value: `0`.
* `gateway_file_share_name` - (Required, ForceNew) The name of the file share. Length from `1` to `255` characters can contain lowercase letters, digits, (.), (_) Or (-), at the same time, must start with a lowercase letter.
* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `ignore_delete` - (Optional) The whether to ignore deleted of the gateway file share. After the opening of the Gateway side delete file or delete cloud (OSS) corresponding to the file. Default value: `false`. **NOTE:** `ignore_delete` and `remote_sync` cannot be enabled simultaneously. Gateway version >= 1.0.40 above support. 
* `in_place` - (Optional, ForceNew) The whether debris optimization of the gateway file share. Default value: `false`.
* `lag_period` - (Optional) The synchronization delay, I.e. gateway local cache sync to Alibaba Cloud Object Storage Service (oss) of the delay time. Unit: `Seconds`. Value range: `5` ~ `120`. Default value: `5`. **NOTE:** Gateway version >= 1.0.40 above support. 
* `local_path` - (Required, ForceNew) The cache disk inside the device name.
* `nfs_v4_optimization` - (Optional) The set up gateway file share NFS protocol, whether to enable NFS v4 optimization improve Mount Upload efficiency. Default value: `false`. **NOTE:** If it is enabled, NFS V3 cannot be mounted. The attribute is valid when the attribute `protocol` is `NFS`. Gateway version >= 1.2.0 above support. 
* `oss_bucket_name` - (Required, ForceNew) The name of the OSS Bucket.
* `oss_bucket_ssl` - (Optional, ForceNew) Whether they are using SSL connect to OSS Bucket.
* `oss_endpoint` - (Required, ForceNew) The gateway file share corresponds to the Object Storage SERVICE (OSS), Bucket Endpoint. **NOTE:** distinguish between intranet and internet Endpoint. We recommend that if the OSS Bucket and the gateway is in the same Region is use the RDS intranet IP Endpoint: `oss-cn-hangzhou-internal.aliyuncs.com`. 
* `partial_sync_paths` - (Optional, ForceNew) In part mode, the directory path group JSON format.
* `path_prefix` - (Optional, ForceNew) The subdirectory path under the object storage (OSS) bucket corresponding to the file share. If it is blank, it means the root directory of the bucket.
* `polling_interval` - (Optional) The reverse synchronization time intervals of the gateway file share. Value range: `15` ~ `36000`. **NOTE:** in copy mode + reverse synchronization is enabled Download file data, value range: `3600` ~ `36000`. 
* `protocol` - (Required, ForceNew) Share types. Valid values: `SMB`, `NFS`.
* `remote_sync` - (Optional) Whether to enable reverse synchronization of the gateway file share. Default value: `false`.
* `remote_sync_download` - (Optional) Copy mode, whether to download the file data. Default value: `false`. **NOTE:** only when the attribute `remote_sync` is `true` or acceded to by the speed synchronization group, this parameter will not take effect. 
* `ro_client_list` - (Optional) File sharing NFS read-only client list (IP address or IP address range). Use commas (,) to separate multiple clients.
* `ro_user_list` - (Optional) The read-only client list. When Protocol for Server Message Block (SMB) to go back to.
* `rw_client_list` - (Optional) Read and write the client list. When Protocol NFS is returned when the status is.
* `rw_user_list` - (Optional) Read-write user list. When Protocol for Server Message Block (SMB) to go back to.
* `squash` - (Optional) The NFS protocol user mapping of the gateway file share. Valid values: `none`, `root_squash`, `all_squash`, `all_anonymous`. Default value: `none`. **NOTE:** The attribute is valid when the attribute `protocol` is `NFS`.
* `support_archive` - (Optional, ForceNew) Whether to support the archive transparent read.
* `transfer_acceleration` - (Optional) The set up gateway file share whether to enable transmission acceleration needs corresponding OSS Bucket enabled transport acceleration. **NOTE:** Gateway version >= 1.3.0 above support. 
* `windows_acl` - (Optional) Whether to enable by Windows access list (requires AD domain) the permissions control. Default value: `false`. **NOTE:** The attribute is valid when the attribute `protocol` is `SMB`. Gateway version >= 1.0.45 above support. 
* `bypass_cache_read` - (Optional) Direct reading OSS of the gateway file share.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway File Share. The value formats as `<gateway_id>:<index_id>`.
* `index_id` - The ID of the file share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway File Share.
* `update` - (Defaults to 5 mins) Used when update the Gateway File Share.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway File Share.

## Import

Cloud Storage Gateway Gateway File Share can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_file_share.example <gateway_id>:<index_id>
```