package opentelekomcloud

import (
	"time"

	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func resourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCTSTrackerCreate,
		Read:   resourceCTSTrackerRead,
		Update: resourceCTSTrackerUpdate,
		Delete: resourceCTSTrackerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_prefix_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_support_smn": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"operations": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"is_send_all_key_operation": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"need_notify_user_list": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}

}

func resourceCTSTrackerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	createOpts := tracker.CreateOpts{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		SimpleMessageNotification: tracker.SimpleMessageNotification{
			IsSupportSMN:          d.Get("is_support_smn").(bool),
			TopicID:               d.Get("topic_id").(string),
			Operations:            resourceCTSOperations(d),
			IsSendAllKeyOperation: d.Get("is_send_all_key_operation").(bool),
			NeedNotifyUserList:    resourceCTSNeedNotifyUserList(d),
		},
	}

	trackers, err := tracker.Create(ctsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating CTS tracker : %s", err)
	}

	d.SetId(trackers.TrackerName)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"enabled"},
		Refresh:    waitForCTSTrackerActive(ctsClient),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for CTS tracker (%s) to become available: %s", trackers.TrackerName, StateErr)
	}

	return resourceCTSTrackerRead(d, meta)

}

func resourceCTSTrackerRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	CtsOpts := tracker.ListOpts{TrackerName: d.Id()}
	trackers, err := tracker.List(ctsClient,CtsOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing cts tracker %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving cts tracker: %s", err)
	}


	ctsTracker := trackers[0]

	d.Set("tracker_name", ctsTracker.TrackerName)
	d.Set("bucket_name", ctsTracker.BucketName)
	d.Set("status", ctsTracker.Status)
	d.Set("file_prefix_name", ctsTracker.FilePrefixName)
	d.Set("is_support_smn", ctsTracker.SimpleMessageNotification.IsSupportSMN)
	d.Set("topic_id", ctsTracker.SimpleMessageNotification.TopicID)
	d.Set("is_send_all_key_operation", ctsTracker.SimpleMessageNotification.IsSendAllKeyOperation)
	d.Set("operations", ctsTracker.SimpleMessageNotification.Operations)
	d.Set("need_notify_user_list", ctsTracker.SimpleMessageNotification.NeedNotifyUserList)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCTSTrackerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	var updateOpts tracker.UpdateOpts

	//as bucket_name is mandatory while updating tracker
	updateOpts.BucketName = d.Get("bucket_name").(string)

	updateOpts.SimpleMessageNotification.TopicID = d.Get("topic_id").(string)

	updateOpts.SimpleMessageNotification.Operations = resourceCTSOperations(d)

	updateOpts.SimpleMessageNotification.NeedNotifyUserList = resourceCTSNeedNotifyUserList(d)

	updateOpts.SimpleMessageNotification.IsSupportSMN = d.Get("is_support_smn").(bool)

	if d.HasChange("file_prefix_name") {
		updateOpts.FilePrefixName = d.Get("file_prefix_name").(string)
	}
	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}
	if d.HasChange("is_send_all_key_operation") {
		updateOpts.SimpleMessageNotification.IsSendAllKeyOperation = d.Get("is_send_all_key_operation").(bool)
	}


	_, err = tracker.Update(ctsClient, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating cts tracker: %s", err)
	}

	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"enabled", "disabled"},
		Target:     []string{"deleted"},
		Refresh:    waitForCTSTrackerDelete(ctsClient),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting cts tracker: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCTSTrackerActive(ctsClient *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := tracker.Get(ctsClient).ExtractTracker()
		if err != nil {
			return nil, "", err
		}

		if n[0].Status == "error" {
			return n, n[0].Status, nil
		}
		return n, n[0].Status, nil
	}
}

func waitForCTSTrackerDelete(ctsClient *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := tracker.Get(ctsClient).ExtractTracker()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted cts tracker")
				return r, "deleted", nil
			}
			return r, "available", err
		}

		tracker_del := tracker.Delete(ctsClient)
		err = tracker_del.ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted cts tracker ")
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "available", nil
				}
			}
			return r, "available", err
		}

		return r, "deleted", nil
	}
}

func resourceCTSOperations(d *schema.ResourceData) []string {
	rawOperations := d.Get("operations").(*schema.Set)
	operation := make([]string, (rawOperations).Len())
	for i, raw := range rawOperations.List() {
		operation[i] = raw.(string)
	}
	return operation
}

func resourceCTSNeedNotifyUserList(d *schema.ResourceData) []string {
	rawNotify := d.Get("need_notify_user_list").(*schema.Set)
	notify := make([]string, (rawNotify).Len())
	for i, raw := range rawNotify.List() {
		notify[i] = raw.(string)
	}
	return notify
}

func flattenStringArray(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func flattenCTSSimpleMessageNotification(trackerObject tracker.Tracker) []map[string]interface{} {
	var smn []map[string]interface{}

	mapping := map[string]interface{}{
		"is_support_smn":            trackerObject.SimpleMessageNotification.IsSupportSMN,
		"topic_id":                  trackerObject.SimpleMessageNotification.TopicID,
		"is_send_all_key_operation": trackerObject.SimpleMessageNotification.IsSendAllKeyOperation,
		"operations":                schema.NewSet(schema.HashString, flattenStringArray(trackerObject.SimpleMessageNotification.Operations)),
		"need_notify_user_list":     schema.NewSet(schema.HashString, flattenStringArray(trackerObject.SimpleMessageNotification.NeedNotifyUserList)),
	}

	smn = append(smn, mapping)

	return smn

}
