package alicloud

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"hash/crc32"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudSsoAccessConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSsoAccessConfigurationCreate,
		Read:   resourceAliCloudCloudSsoAccessConfigurationRead,
		Update: resourceAliCloudCloudSsoAccessConfigurationUpdate,
		Delete: resourceAliCloudCloudSsoAccessConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-zA-z0-9-]{1,32}$`), "The name of the resource. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-)"),
			},
			"session_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(900, 43200),
			},
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(0, 1024),
			},
			"force_remove_permission_policies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"permission_policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					policy := v.(map[string]interface{})
					if v, ok := policy["permission_policy_type"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := policy["permission_policy_name"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := policy["permission_policy_document"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					return int(crc32.ChecksumIEEE([]byte(buf.String())))
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					log.Println("policies中的原k：", k)
					log.Println("policies中的原old：", old)
					log.Println("policies中的原new：", new)
					log.Println("policies中的原d.Get(k)：", d.Get(k))
					log.Println("policies中的原old != new：", old != new)
					log.Println("policies中的原old != d.Get(k)：", old != d.Get(k))
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_policy_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"System", "Inline"}, false),
						},
						"permission_policy_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permission_policy_document": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.ValidateJsonString,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								log.Println("document中的原k：", k)
								log.Println("document中的原old：", old)
								log.Println("document中的原new：", new)
								log.Println("document中的原d.Get(k)：", d.Get(k))
								log.Println("document中的原old != new：", old != new)
								log.Println("document中的原old != d.Get(k)：", old != d.Get(k))
								//if old != "" && new != "" && old != new {
								//	log.Println("原old：", old)
								//	log.Println("原new：", new)
								//	log.Println("原old != new：", old != new)
								//	old = Trim(strings.Replace(old, "\n", "", -1))
								//	new = Trim(strings.Replace(new, "\n", "", -1))
								//	log.Println("优化old：", old)
								//	log.Println("优化new：", new)
								//	log.Println("优化old != new：", old != new)
								//	equal, _ := compareJsonTemplateAreEquivalent(old, new)
								//	log.Println("equal：", equal)
								//}
								return false
							},
						},
					},
				},
			},
			"access_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudSsoAccessConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessConfiguration"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	request["DirectoryId"] = d.Get("directory_id")
	request["AccessConfigurationName"] = d.Get("access_configuration_name")

	if v, ok := d.GetOk("session_duration"); ok {
		request["SessionDuration"] = v
	}

	if v, ok := d.GetOk("relay_state"); ok {
		request["RelayState"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_access_configuration", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.AccessConfiguration", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_sso_access_configuration")
	} else {
		accessConfigurationId := resp.(map[string]interface{})["AccessConfigurationId"]
		d.SetId(fmt.Sprintf("%v:%v", request["DirectoryId"], accessConfigurationId))
	}

	return resourceAliCloudCloudSsoAccessConfigurationUpdate(d, meta)
}

func resourceAliCloudCloudSsoAccessConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}

	object, err := cloudssoService.DescribeCloudSsoAccessConfiguration(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_configuration cloudssoService.DescribeCloudSsoAccessConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("directory_id", parts[0])
	d.Set("access_configuration_name", object["AccessConfigurationName"])
	d.Set("relay_state", object["RelayState"])
	d.Set("description", object["Description"])
	d.Set("access_configuration_id", object["AccessConfigurationId"])

	if sessionDuration, ok := object["SessionDuration"]; ok && fmt.Sprint(sessionDuration) != "0" {
		d.Set("session_duration", formatInt(sessionDuration))
	}

	listPermissionPoliciesInAccessConfigurationObject, err := cloudssoService.ListPermissionPoliciesInAccessConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if permissionPolicies, ok := listPermissionPoliciesInAccessConfigurationObject["PermissionPolicies"]; ok && permissionPolicies != nil {
		permissionPoliciesMaps := make([]map[string]interface{}, 0)
		for _, permissionPoliciesList := range permissionPolicies.([]interface{}) {
			permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})
			permissionPoliciesMap := make(map[string]interface{})

			if permissionPolicyType, ok := permissionPoliciesArg["PermissionPolicyType"]; ok {
				permissionPoliciesMap["permission_policy_type"] = permissionPolicyType
			}

			if permissionPolicyName, ok := permissionPoliciesArg["PermissionPolicyName"]; ok {
				permissionPoliciesMap["permission_policy_name"] = permissionPolicyName
			}

			if permissionPolicyDocument, ok := permissionPoliciesArg["PermissionPolicyDocument"]; ok {
				permissionPoliciesMap["permission_policy_document"] = permissionPolicyDocument
			}

			permissionPoliciesMaps = append(permissionPoliciesMaps, permissionPoliciesMap)
		}

		d.Set("permission_policies", permissionPoliciesMaps)
	}

	return nil
}

func resourceAliCloudCloudSsoAccessConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	updateAccessConfigurationReq := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("session_duration") {
		update = true

		if v, ok := d.GetOk("session_duration"); ok {
			updateAccessConfigurationReq["NewSessionDuration"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("relay_state") {
		update = true
	}
	if v, ok := d.GetOk("relay_state"); ok {
		updateAccessConfigurationReq["NewRelayState"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		updateAccessConfigurationReq["NewDescription"] = v
	}

	if update {
		action := "UpdateAccessConfiguration"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, updateAccessConfigurationReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAccessConfigurationReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("session_duration")
		d.SetPartial("relay_state")
		d.SetPartial("description")
	}

	if d.HasChange("permission_policies") {
		oldPermissionPolicies, newPermissionPolicies := d.GetChange("permission_policies")
		removed := oldPermissionPolicies.(*schema.Set).Difference(newPermissionPolicies.(*schema.Set)).List()
		added := newPermissionPolicies.(*schema.Set).Difference(oldPermissionPolicies.(*schema.Set)).List()

		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}

		if len(removed) > 0 {
			action := "RemovePermissionPolicyFromAccessConfiguration"

			removePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range removed {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, removePermissionPolicyFromAccessConfigurationReq, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removePermissionPolicyFromAccessConfigurationReq)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(added) > 0 {
			action := "AddPermissionPolicyToAccessConfiguration"

			addPermissionPolicyToAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range added {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				addPermissionPolicyToAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				addPermissionPolicyToAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				if addPermissionPolicyToAccessConfigurationReq["PermissionPolicyType"] == "Inline" {

					if inlinePolicyDocument, ok := permissionPoliciesArg["permission_policy_document"]; ok {
						addPermissionPolicyToAccessConfigurationReq["InlinePolicyDocument"] = inlinePolicyDocument
					}
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, addPermissionPolicyToAccessConfigurationReq, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, addPermissionPolicyToAccessConfigurationReq)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		// Provisioning access configuration when permission policies has changed.
		objects, err := cloudssoService.DescribeCloudSsoAccessConfigurationProvisionings(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]))
		if err != nil {
			return WrapError(err)
		}

		for _, object := range objects {
			err = cloudssoService.CloudssoServicAccessConfigurationProvisioning(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]), fmt.Sprint(object["TargetType"]), fmt.Sprint(object["TargetId"]))
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("permission_policies")
	}

	d.Partial(false)

	return resourceAliCloudCloudSsoAccessConfigurationRead(d, meta)
}

func resourceAliCloudCloudSsoAccessConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("permission_policies"); ok {
		removed := v.(*schema.Set).List()

		if len(removed) > 0 {
			action := "RemovePermissionPolicyFromAccessConfiguration"

			removePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range removed {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, removePermissionPolicyFromAccessConfigurationReq, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removePermissionPolicyFromAccessConfigurationReq)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	action := "DeleteAccessConfiguration"

	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	if v, ok := d.GetOkExists("force_remove_permission_policies"); ok {
		request["ForceRemovePermissionPolicies"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"DeletionConflict.AccessConfiguration.Provisioning", "DeletionConflict.AccessConfiguration.AccessAssignment", "OperationConflict.Task", "DeletionConflict.AccessConfiguration.Task"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
