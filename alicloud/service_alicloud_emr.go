package alicloud

import (
	"fmt"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type EmrService struct {
	client *connectivity.AliyunClient
}

func (s *EmrService) DescribeEmrCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeClusterV2"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Id":       id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	if v.(map[string]interface{})["Status"] == "RELEASED" {
		return object, WrapErrorf(NotFoundErr("EmrCluster", id), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EmrService) GetEmrV2Cluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetCluster"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Cluster", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Cluster", response)
	}
	if v.(map[string]interface{})["ClusterState"] == "TERMINATED" {
		return object, WrapErrorf(NotFoundErr("EmrCluster", id), NotFoundWithResponse, response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EmrService) GetEmrV2Operation(id string, operationID string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetOperation"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ClusterId":   id,
		"OperationId": operationID,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Operation", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Operation", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EmrService) ListEmrV2NodeGroups(clusterId string, nodeGroupIds []string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListNodeGroups"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClusterId":       clusterId,
		"NodeGroupStates": []string{"RUNNING"},
		"MaxResults":      PageSizeXLarge,
	}
	if nodeGroupIds != nil && len(nodeGroupIds) > 0 {
		request["NodeGroupIds"] = nodeGroupIds
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, clusterId, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NodeGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, clusterId, "$.NodeGroups", response)
	}
	if v != nil {
		object = v.([]interface{})
	}
	return object, nil
}

func (s *EmrService) DataSourceDescribeEmrCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeClusterV2"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Id":       id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EmrService) WaitForEmrCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object["Id"].(string) == id && status != Deleted {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["Id"].(string), id, ProviderERROR)
		}
	}
	return nil
}

func (s *EmrService) WaitForEmrV2Cluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.GetEmrV2Cluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object["ClusterId"].(string) == id && status != Deleted {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ClusterId"].(string), id, ProviderERROR)
		}
	}
	return nil
}

func (s *EmrService) WaitForEmrV2Operation(id string, nodeGroupId string, operationID string, timeout int, wg *sync.WaitGroup, cm *sync.Map) {
	defer wg.Done()
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, _ := s.GetEmrV2Operation(id, operationID)
		if object != nil {
			if object["OperationState"].(string) == "COMPLETED" || object["OperationState"].(string) == "PARTIAL_COMPLETED" {
				cm.Store(nodeGroupId, true)
				return
			}

			if object["OperationState"].(string) == "FAILED" || object["OperationState"].(string) == "TERMINATED" {
				cm.Store(nodeGroupId, false)
				return
			}
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
		if time.Now().After(deadline) {
			cm.Store(nodeGroupId, false)
			return
		}
	}
}

func (s *EmrService) EmrClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}

		return object, object["Status"].(string), nil
	}
}

func (s *EmrService) EmrV2ClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetEmrV2Cluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ClusterState"].(string) == failState {
				return object, object["ClusterState"].(string), WrapError(Error(FailedToReachTargetStatus, object["ClusterState"].(string)))
			}
		}
		return object, object["ClusterState"].(string), nil
	}
}

func (s *EmrService) setEmrClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := emr.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(TagResourceCluster)
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := emr.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(TagResourceCluster)
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *EmrService) DescribeEmrClusterTags(resourceId string, resourceType TagResourceType) (tags []emr.TagResource, err error) {
	request := emr.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithEmrClient(func(client *emr.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*emr.ListTagResourcesResponse)
	tags = response.TagResources.TagResource

	return
}

func (s *EmrService) diffTags(oldTags, newTags []emr.TagResourcesTag) ([]emr.TagResourcesTag, []emr.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []emr.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *EmrService) tagsFromMap(m map[string]interface{}) []emr.TagResourcesTag {
	result := make([]emr.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, emr.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *EmrService) tagsToMap(tags []emr.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.TagKey] = t.TagValue
	}
	return result
}

func (s *EmrService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Emr", "2016-04-08", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *EmrService) SetEmrClusterTagsNew(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		client := s.client
		_, nraw := d.GetChange("tags")

		var createTags []map[string]interface{}
		newTagMap := nraw.(map[string]interface{})
		for newKey, newValue := range newTagMap {
			createTags = append(createTags, map[string]interface{}{
				"Key":   newKey,
				"Value": newValue,
			})
		}

		tags, err := s.ListTagResourcesNew(d.Id(), string(TagResourceCluster))
		if err != nil {
			return WrapError(err)
		}

		var deleteTagKeys []string
		for oldKey, oldValue := range tagsToMap(tags) {
			newValue, ok := newTagMap[oldKey]
			if !ok || oldValue != newValue {
				deleteTagKeys = append(deleteTagKeys, oldKey)
			}
		}

		if len(createTags) > 0 {
			action := "TagResources"
			tagResourcesRequest := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": string(TagResourceCluster),
				"ResourceIds":  []string{d.Id()},
				"Tags":         createTags,
			}
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = client.RpcPost("Emr", "2021-03-20", action, nil, tagResourcesRequest, false)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				addDebug(action, d.Id(), tagResourcesRequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(deleteTagKeys) > 0 {
			action := "UntagResources"
			unTagResourcesRequest := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": string(TagResourceCluster),
				"ResourceIds":  []string{d.Id()},
				"Tags":         deleteTagKeys,
			}
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = client.RpcPost("Emr", "2021-03-20", action, nil, unTagResourcesRequest, false)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				addDebug(action, d.Id(), unTagResourcesRequest)
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

func (s *EmrService) ListTagResourcesNew(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceIds":  []string{id},
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Emr", "2021-03-20", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *EmrService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

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
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Emr", "2016-04-08", action, nil, request, false)
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
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Emr", "2016-04-08", action, nil, request, false)
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
func (s *EmrService) DescribeClusterBasicInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeClusterBasicInfo"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *EmrService) DescribeClusterV2(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeClusterV2"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Id":       id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EmrService) DescribeEmrMainVersionClusterTypes(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeEmrMainVersion"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EmrVersion": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.EmrMainVersion.ClusterTypeInfoList.ClusterTypeInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.EmrMainVersion.ClusterTypeInfoList.ClusterTypeInfo", response)
	}
	object = v.([]interface{})
	return object, nil
}
