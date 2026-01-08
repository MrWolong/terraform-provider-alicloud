package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall Instance. >>> Resource test cases, automatically generated.
// Case 国内版按量付费2.0 11709
func TestAccAliCloudCloudFirewallInstanceV2_basic11709(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11709_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_cn",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_cn",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

var AliCloudCloudFirewallInstanceV2Map11709 = map[string]string{
	"create_time": CHECKSET,
	"end_time":    CHECKSET,
	"user_status": CHECKSET,
	"status":      CHECKSET,
}

func AliCloudCloudFirewallInstanceV2BasicDependence11709(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 国际版按量付费2.0 11710
func TestAccAliCloudCloudFirewallInstanceV2_basic11710(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sdl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "false",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstanceV2_basic11710_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallInstanceV2Map11709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallInstanceV2BasicDependence11709)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"product_code": "cfw",
					"product_type": "cfw_elasticity_public_intl",
					"spec":         "payg_version",
					"sdl":          "true",
					"cfw_log":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"product_code": "cfw",
						"product_type": "cfw_elasticity_public_intl",
						"spec":         "payg_version",
						"sdl":          "true",
						"cfw_log":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

// Test CloudFirewall Instance. <<< Resource test cases, automatically generated.
