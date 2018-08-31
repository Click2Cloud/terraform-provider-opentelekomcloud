package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func resourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSBSBackupPolicyCreate,
		Read:   resourceCSBSBackupPolicyRead,
		Update: resourceCSBSBackupPolicyUpdate,
		Delete: resourceCSBSBackupPolicyDelete,
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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//ForceNew:     false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:     false,
			},
			"provider_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"common": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"scheduled_operations": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheduled_period_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew:     false,
						},
						"scheduled_period_description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: false,
						},
						"max_backups": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"retention_duration_days": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"permanent": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"plan_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"pattern": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"operation_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"scheduled_period_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"scheduler_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resources": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}

}

func resourceCSBSBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup policy Client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Policy: policies.PolicyCreate{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ProviderId:  "fc4d5750-22e7-4798-8a46-f48f62c4c1da",
			Parameters: policies.PolicyParam{
				Common: map[string]interface{}{},
			},
			ScheduledOperations: resourceCSBSScheduleV1(d),

			Resources: resourceCSBSResourceV1(d),
			Tags:      resourceCSBSPolicyTagsV1(d),
		},
	}

	create, err := policies.Create(policyClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy : %s", err)
	}
	/*if err != nil {
		if _, ok := err.(gophercloud.ErrDefault500); ok {
			return fmt.Errorf("Server (%s) is already in service.", createOpts.Policy.Resources[0].Id)
		}
		return fmt.Errorf("Error creating OpenTelekomCloud Backup : %s", err)
	} else {
		log.Printf("[DEBUG] Create Id (%s): %s", create.Id)
	}*/

	d.SetId(create.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"suspended"},
		Refresh:    waitForCSBSBackupPolicyActive(policyClient, create.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for Backup Policy (%s) to become available: %s", create.Id, StateErr)
	}

	return resourceCSBSBackupPolicyRead(d, meta)

}

func resourceCSBSBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup policy: %s", err)
	}

	n, err := policies.Get(policyClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud backup policy: %s", err)
	}
	var resourcelist []map[string]interface{}
	for _, resources := range n.Resources {
		mapping := map[string]interface{}{
			"resource_id":   resources.Id,
			"resource_type": resources.Type,
			"resource_name": resources.Name,
		}
		resourcelist = append(resourcelist, mapping)
	}

	var scheduledlist []map[string]interface{}
	for _, schedule := range n.ScheduledOperations {
		mapping := map[string]interface{}{
			"scheduled_period_description": schedule.Description,
			"enabled":                      schedule.Enabled,
			"trigger_id":                   schedule.TriggerId,
			"scheduled_period_name":        schedule.Name,
			"operation_type":               schedule.OperationType,
			"max_backups":                  schedule.OperationDefinition.MaxBackups,
			"retention_duration_days":      schedule.OperationDefinition.RetentionDurationDays,
			"permanent":                    schedule.OperationDefinition.Permanent,
			"plan_id":                      schedule.OperationDefinition.PlanId,
			"scheduler_id":                 schedule.Trigger.Id,
			"scheduler_name":               schedule.Trigger.Name,
			"scheduler_type":               schedule.Trigger.Type,
			"pattern":                      schedule.Trigger.Properties.Pattern,
			"scheduled_period_id":          schedule.Id,
		}
		scheduledlist = append(scheduledlist, mapping)
	}

	var tagslist []map[string]interface{}
	for _, tag := range n.Tags {
		mapping := map[string]interface{}{
			"key":   tag.Key,
			"value": tag.Value,
		}
		tagslist = append(tagslist, mapping)
	}

	d.Set("description", n.Description)
	d.Set("id", n.Id)
	d.Set("name", n.Name)
	d.Set("common", n.Parameters.Common)
	d.Set("project_id", n.ProjectId)
	d.Set("provider_id", n.ProviderId)
	d.Set("status", n.Status)
	if err := d.Set("resources", resourcelist); err != nil {
		return err
	}
	if err := d.Set("scheduled_operations", scheduledlist); err != nil {
		return err
	}
	if err := d.Set("tags", tagslist); err != nil {
		return err
	}
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCSBSBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Backup Policy: %s", err)
	}
	var updateOpts policies.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Policy.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Policy.Description = d.Get("description").(string)
	}
	updateOpts.Policy.Parameters.Common = map[string]interface{}{}

	if d.HasChange("resources") {
		updateOpts.Policy.Resources = resourceCSBSResourceUpdateV1(d)
	}
	if d.HasChange("scheduled_operations") {
		updateOpts.Policy.ScheduledOperations = resourceCSBScheduleUpdateV1(d)
	}

	_, err = policies.Update(policyClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud RTS Software Deployment: %s", err)
	}

	return resourceCSBSBackupPolicyRead(d, meta)
}

func resourceCSBSBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup Policy: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"deleted"},
		Refresh:    waitForVBSPolicyDelete(policyClient, d.Id()),
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

func waitForCSBSBackupPolicyActive(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := policies.Get(policyClient, policyID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "error" {
			return n, n.Status, nil
		}
		return n, n.Status, nil
	}
}

func waitForVBSPolicyDelete(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := policies.Get(policyClient, policyID).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Backup Policy %s", policyID)
				return r, "deleted", nil
			}
			return r, "available", err
		}

		policy := policies.Delete(policyClient, policyID)
		err = policy.Err
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Backup Policy %s", policyID)
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

func resourceCSBSScheduleV1(d *schema.ResourceData) []policies.ScheduledOperations {
	rawTags := d.Get("scheduled_operations").([]interface{})
	schedule := make([]policies.ScheduledOperations, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		schedule[i] = policies.ScheduledOperations{
			Name:          rawMap["scheduled_period_name"].(string),
			Description:   rawMap["scheduled_period_description"].(string),
			TriggerId:     rawMap["trigger_id"].(string),
			Enabled:       rawMap["enabled"].(bool),
			OperationType: rawMap["operation_type"].(string),
			Trigger: policies.Trigger{
				Properties: policies.TriggerProperties{
					Pattern: rawMap["pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinition{
				MaxBackups:            rawMap["max_backups"].(string),
				RetentionDurationDays: rawMap["retention_duration_days"].(string),
				Permanent:             rawMap["permanent"].(string),
				PlanId:                rawMap["plan_id"].(string),
			},
		}
	}
	return schedule
}

func resourceCSBSResourceV1(d *schema.ResourceData) []policies.Resource {
	rawTags := d.Get("resources").([]interface{})
	resources := make([]policies.Resource, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		resources[i] = policies.Resource{
			Name: rawMap["resource_name"].(string),
			Id:   rawMap["resource_id"].(string),
			Type: rawMap["resource_type"].(string),
		}
	}
	return resources
}

func resourceCSBSPolicyTagsV1(d *schema.ResourceData) []policies.ResourceTag {
	rawTags := d.Get("tags").([]interface{})
	tags := make([]policies.ResourceTag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.ResourceTag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
}

func resourceCSBScheduleUpdateV1(d *schema.ResourceData) []policies.ScheduledOperationsUpdate {
	rawTags := d.Get("scheduled_operations").([]interface{})
	schedule := make([]policies.ScheduledOperationsUpdate, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		schedule[i] = policies.ScheduledOperationsUpdate{
			Id:          rawMap["scheduled_period_id"].(string),
			Name:        rawMap["scheduled_period_name"].(string),
			Description: rawMap["scheduled_period_description"].(string),
			Enabled:     rawMap["enabled"].(bool),
			Trigger: policies.TriggerUpdate{
				Properties: policies.TriggerPropertiesUpdate{
					Pattern: rawMap["pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinitionUpdate{
				MaxBackups:            rawMap["max_backups"].(string),
				RetentionDurationDays: rawMap["retention_duration_days"].(string),
				Permanent:             rawMap["permanent"].(string),
			},
		}
	}
	return schedule
}

func resourceCSBSResourceUpdateV1(d *schema.ResourceData) []policies.ResourceUpdate {
	rawTags := d.Get("resources").([]interface{})
	resources := make([]policies.ResourceUpdate, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		resources[i] = policies.ResourceUpdate{
			Name: rawMap["resource_name"].(string),
			Id:   rawMap["resource_id"].(string),
			Type: rawMap["resource_type"].(string),
		}
	}
	return resources
}
