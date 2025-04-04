package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EcdService struct {
	client *connectivity.AliyunClient
}

func (s *EcdService) DescribeEcdPolicyGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribePolicyGroups"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"PolicyGroupId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DescribePolicyGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DescribePolicyGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PolicyGroupId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func parsingResourceType(d *schema.ResourceData, resourceType string) ([]interface{}, []interface{}) {
	oraw, nraw := d.GetChange(resourceType)
	removed := oraw.(*schema.Set).List()
	added := nraw.(*schema.Set).List()
	return removed, added
}

func (s *EcdService) setAuthAccessPolicyRules(d *schema.ResourceData, request map[string]interface{}, resourceType string) error {
	if d.HasChange(resourceType) {
		removed, added := parsingResourceType(d, resourceType)
		client := s.client
		var err error
		if len(removed) > 0 {
			action := "ModifyPolicyGroup"
			req := map[string]interface{}{
				"PolicyGroupId": d.Id(),
			}
			var response map[string]interface{}
			for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range removed {
				authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
				req["RevokeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
				req["RevokeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("ecd", "2020-09-30", action, nil, req, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidPolicyStatus.Modification"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range added {
				authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
				request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
				request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
			}
		}
	}
	return nil
}

func (s *EcdService) setAuthSecurityPolicyRules(d *schema.ResourceData, request map[string]interface{}, resourceType string) error {
	if d.HasChange(resourceType) {
		removed, added := parsingResourceType(d, resourceType)
		var err error
		client := s.client
		if len(removed) > 0 {
			action := "ModifyPolicyGroup"
			req := map[string]interface{}{
				"PolicyGroupId": d.Id(),
			}
			var response map[string]interface{}
			for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range removed {
				authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("ecd", "2020-09-30", action, nil, req, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidPolicyStatus.Modification"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range added {
				authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
			}
		}
	}
	return nil
}

func (s *EcdService) DescribeEcdSimpleOfficeSite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeOfficeSites"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"OfficeSiteId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.OfficeSites", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.OfficeSites", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["OfficeSiteId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdSimpleOfficeSiteStateRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdSimpleOfficeSite(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		return object, object["Status"].(string), nil
	}
}

func (s *EcdService) DescribeEcdNasFileSystem(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeNASFileSystems"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"FileSystemId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.FileSystems", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FileSystems", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["FileSystemId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) DescribeEcdNetworkPackage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeNetworkPackages"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"NetworkPackageId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.NetworkPackages", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkPackages", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["NetworkPackageId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdNasFileSystemStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdNasFileSystem(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["FileSystemStatus"]) == failState {
				return object, fmt.Sprint(object["FileSystemStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["FileSystemStatus"])))
			}
		}
		return object, fmt.Sprint(object["FileSystemStatus"]), nil
	}
}

func (s *EcdService) EcdNetworkPackageRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdNetworkPackage(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["NetworkPackageStatus"]) == failState {
				return object, fmt.Sprint(object["NetworkPackageStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["FileSystemStatus"])))
			}
		}
		return object, fmt.Sprint(object["NetworkPackageStatus"]), nil
	}
}

func (s *EcdService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId":   []string{d.Id()},
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId":   []string{d.Id()},
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *EcdService) DescribeEcdDesktop(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDesktops"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DesktopId.1": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Desktops", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Desktops", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DesktopId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdDesktopStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdDesktop(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DesktopStatus"]) == failState {
				return object, fmt.Sprint(object["DesktopStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DesktopStatus"])))
			}
		}

		return object, fmt.Sprint(object["DesktopStatus"]), nil
	}
}

func (s *EcdService) EcdDesktopDesktopTypeRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdDesktop(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DesktopType"]) == failState {
				return object, fmt.Sprint(object["DesktopType"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DesktopType"])))
			}
		}

		return object, fmt.Sprint(object["DesktopType"]), nil
	}
}

func (s *EcdService) EcdDesktopChargeTypeFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdDesktop(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(convertEcdDesktopPaymentTypeResponse(object["ChargeType"])) == failState {
				return object, fmt.Sprint(object["ChargeType"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(convertEcdDesktopPaymentTypeResponse(object["DesktopStatus"]))))
			}
		}

		return object, fmt.Sprint(convertEcdDesktopPaymentTypeResponse(object["ChargeType"])), nil
	}
}

func (s *EcdService) DescribeEcdImage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeImages"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ImageId":  []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Images", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Images", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ImageId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdImageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdImage(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *EcdService) DescribeEcdCommand(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client

	action := "DescribeInvocations"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"InvokeId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Invocations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Invocations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {

		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InvokeId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdCommandStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdCommand(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["InvocationStatus"]) == failState {
				return object, fmt.Sprint(object["InvocationStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InvocationStatus"])))
			}
		}
		return object, fmt.Sprint(object["InvocationStatus"]), nil
	}
}

func (s *EcdService) DescribeEcdSnapshot(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"SnapshotId": id,
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeSnapshots"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Snapshots", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Snapshots", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Snapshot", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *EcdService) DescribeEcdBundle(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeBundles"

	client := s.client

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"BundleId": []string{id},
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Bundles", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Bundles", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD:Bundle", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["BundleId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("ECD:Bundle", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *EcdService) DescribeEcdRamDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDirectories"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DirectoryId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directories", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Directories", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DirectoryId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) DescribeEcdAdConnectorDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDirectories"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DirectoryId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directories", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Directories", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DirectoryId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) DescribeEcdAdConnectorOfficeSite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeOfficeSites"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"OfficeSiteId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.OfficeSites", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.OfficeSites", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["OfficeSiteId"]) != id {
			return object, WrapErrorf(NotFoundErr("ECD", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EcdService) EcdAdConnectorOfficeSiteStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdAdConnectorOfficeSite(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
