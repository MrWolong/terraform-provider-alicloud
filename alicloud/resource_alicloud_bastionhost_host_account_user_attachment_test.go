package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudBastionhostHostAccountUserAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_account_user_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudBastionhostHostAccountUserAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAccountUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudBastionhostHostAccountUserAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":          "${alicloud_bastionhost_user.default.user_id}",
					"host_id":          "${alicloud_bastionhost_host_account.default.0.host_id}",
					"instance_id":      "${alicloud_bastionhost_host_account.default.0.instance_id}",
					"host_account_ids": []string{"${alicloud_bastionhost_host_account.default.0.host_account_id}", "${alicloud_bastionhost_host_account.default.1.host_account_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":            CHECKSET,
						"host_id":            CHECKSET,
						"instance_id":        CHECKSET,
						"host_account_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_ids": []string{"${alicloud_bastionhost_host_account.default.0.host_account_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_ids": []string{"${alicloud_bastionhost_host_account.default.0.host_account_id}", "${alicloud_bastionhost_host_account.default.1.host_account_id}", "${alicloud_bastionhost_host_account.default.2.host_account_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_ids.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudBastionhostHostAccountUserAttachmentMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AliCloudBastionhostHostAccountUserAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_bastionhost_instances" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		count  = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? 0 : 1
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_bastionhost_instance" "default" {
  		count              = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? 0 : 1
  		description        = var.name
  		license_code       = "bhah_ent_50_asset"
  		plan_code          = "cloudbastion"
  		storage            = "5"
  		bandwidth          = "5"
  		period             = "1"
  		vswitch_id         = data.alicloud_vswitches.default.ids.0
  		security_group_ids = [alicloud_security_group.default.0.id]
	}

	resource "alicloud_bastionhost_host" "default" {
  		instance_id          = local.instance_id
  		host_name            = var.name
  		active_address_type  = "Private"
  		host_private_address = "172.16.0.10"
  		os_type              = "Linux"
  		source               = "Local"
	}

	resource "alicloud_bastionhost_host_account" "default" {
  		count             = 3
  		instance_id       = alicloud_bastionhost_host.default.instance_id
  		host_account_name = "${var.name}-${count.index}"
  		host_id           = alicloud_bastionhost_host.default.host_id
  		protocol_name     = "SSH"
  		password          = "YourPassword12345"
	}

	resource "alicloud_bastionhost_user" "default" {
  		instance_id         = alicloud_bastionhost_host.default.instance_id
  		mobile              = "13312345678"
  		mobile_country_code = "CN"
  		password            = "YourPassword-123"
  		source              = "Local"
  		user_name           = var.name
	}

	locals {
  		instance_id = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? data.alicloud_bastionhost_instances.default.ids.0 : alicloud_bastionhost_instance.default.0.id
	}
`, name)
}
