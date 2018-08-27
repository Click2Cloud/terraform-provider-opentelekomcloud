package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/policies"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/tags"
	"log"
	"time"
)

func resourceVBSBackupPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVBSBackupPolicyV2Create,
		Read:   resourceVBSBackupPolicyV2Read,
		Update: resourceVBSBackupPolicyV2Update,
		Delete: resourceVBSBackupPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"backup_policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"start_time": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"rentention_num": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"retain_first_backup": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"policy_resource_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVBSBackupPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy Client: %s", err)
	}

	createOpts := policies.CreateOpts{
		PolicyName: d.Get("backup_policy_name").(string),
		ScheduledPolicy: policies.CreateSchedule{
			StartTime:         d.Get("start_time").(string),
			Frequency:         d.Get("frequency").(int),
			RententionNum:     d.Get("rentention_num").(int),
			RemainFirstBackup: d.Get("retain_first_backup").(string),
			Status:            d.Get("status").(string),
		},
		Tags: resourceVBSTagsV2(d),
	}

	create, err := policies.Create(vbsClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy: %s", err)
	}
	d.SetId(create.PolicyID)

	log.Printf("[DEBUG] Waiting for OpenTelekomcomCloud Backup Policy (%s) to become available", d.Id())

	stateConf := &resource.StateChangeConf{
		//	Pending:    []string{"creating"},
		Target:     []string{"ON", "OFF"},
		Refresh:    waitForVBSPolicyActive(vbsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, Stateerr := stateConf.WaitForState()
	if Stateerr != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy : %s", Stateerr)
	}

	return resourceVBSBackupPolicyV2Read(d, meta)

}

func resourceVBSBackupPolicyV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy Client: %s", err)
	}

	PolicyOpts := policies.ListOpts{PolicyID: d.Id()} //shares.Get(vbsV2Client, d.Id()).Extract()
	policies, err := policies.List(vbsClient, PolicyOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Backup Policy: %s", err)
	}

	n := policies[0]

	d.Set("backup_policy_name", n.PolicyName)
	d.Set("start_time", n.ScheduledPolicy.StartTime)
	d.Set("frequency", n.ScheduledPolicy.Frequency)
	d.Set("rentention_num", n.ScheduledPolicy.RententionNum)
	d.Set("retain_first_backup", n.ScheduledPolicy.RemainFirstBackup)
	d.Set("status", n.ScheduledPolicy.Status)

	tags, err := tags.Get(vbsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}
		return fmt.Errorf("Error retrieving OpenTelekomCloud Backup Policy Tags: %s", err)
	}
	var tagList []map[string]interface{}
	for _, v := range tags.Tag {
		tag := make(map[string]interface{})
		tag["key"] = v.Key
		tag["value"] = v.Value

		tagList = append(tagList, tag)
	}
	if err := d.Set("tags", tagList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tags to state for OpenTelekomCloud backup policy (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceVBSBackupPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Share File: %s", err)
	}
	var updateOpts policies.UpdateOpts

	if d.HasChange("backup_policy_name") || d.HasChange("start_time") || d.HasChange("frequency") ||
		d.HasChange("rentention_num") || d.HasChange("retain_first_backup") || d.HasChange("status") {
		if d.HasChange("backup_policy_name") {
			updateOpts.PolicyName = d.Get("backup_policy_name").(string)
		}
		if d.HasChange("start_time") {
			updateOpts.ScheduledPolicy.StartTime = d.Get("start_time").(string)
		}
		if d.HasChange("frequency") {
			updateOpts.ScheduledPolicy.Frequency = d.Get("frequency").(int)
		}
		if d.HasChange("rentention_num") {
			updateOpts.ScheduledPolicy.RententionNum = d.Get("rentention_num").(int)
		}
		if d.HasChange("retain_first_backup") {
			updateOpts.ScheduledPolicy.RemainFirstBackup = d.Get("retain_first_backup").(string)
		}
		if d.HasChange("status") {
			updateOpts.ScheduledPolicy.Status = d.Get("status").(string)
		}
		_, err = policies.Update(vbsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud backup policy: %s", err)
		}
	}
	if d.HasChange("tags") {
		//on tags update , delete the old ones and add new tags.
		// first get the list of old tags
		oldTags, _ := tags.Get(vbsClient, d.Id()).Extract()
		//delete the old tags
		deleteopts := tags.BatchOpts{Action: "delete", Tags: oldTags.Tag}
		delete := tags.BatchAction(vbsClient, d.Id(), deleteopts)
		if delete.Err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud backup policy tags: %s", delete.Err)
		}
		//add the new new tags
		createTags := tags.BatchAction(vbsClient, d.Id(), tags.BatchOpts{Action: "create", Tags: resourceVBSUpdateTagsV2(d)})
		if createTags.Err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud backup policy tags: %s", createTags.Err)
		}
	}
	return resourceVBSBackupPolicyV2Read(d, meta)
}

func resourceVBSBackupPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"deleted"},
		Refresh:    waitForVBSPolicyDelete(vbsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Backup Policy: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVBSPolicyDelete(vbsClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := policies.List(vbsClient, policies.ListOpts{PolicyID: policyID})

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", policyID)
				return r, "deleted", nil
			}
			return r, "available", err
		}
		delete := policies.Delete(vbsClient, policyID)
		err = delete.Err
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", policyID)
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

func waitForVBSPolicyActive(vbsClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		PolicyOpts := policies.ListOpts{PolicyID: policyID} //shares.Get(vbsV2Client, d.Id()).Extract()
		policies, err := policies.List(vbsClient, PolicyOpts)
		if err != nil {
			return nil, "", err
		}
		n := policies[0]
		//if n.ScheduledPolicy.Status != "ON" || n.ScheduledPolicy.Status != "OFF" {
		//	return n, "creating", nil
		//}

		return n, n.ScheduledPolicy.Status, nil
	}
}

func resourceVBSTagsV2(d *schema.ResourceData) []policies.Tags {
	rawTags := d.Get("tags").([]interface{})
	tags := make([]policies.Tags, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.Tags{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
}

func resourceVBSUpdateTagsV2(d *schema.ResourceData) []tags.ActionTags {
	rawTags := d.Get("tags").([]interface{})
	tagList := make([]tags.ActionTags, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tagList[i] = tags.ActionTags{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tagList
}
