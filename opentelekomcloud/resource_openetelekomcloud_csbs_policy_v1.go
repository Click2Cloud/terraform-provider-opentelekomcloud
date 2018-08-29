package opentelekomcloud
import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"log"
	"time"
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
						"so_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew:     false,
						},
						"so_description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew:     false,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default: true,
							ForceNew:     false,
						},
						"max_backups": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew:     false,
						},
						"pattern": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew:     false,
						},
						"operation_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew:false,
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
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"r_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

					},
				},
			},
			/*"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew:true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew:true,
						},
					},
				},
			},*/
		},
	}

	}


func resourceCSBSBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup policy Client: %s", err)
	}
	log.Printf("[INFO] create opts ")

	createOpts:=policies.CreateOpts{
		Policy: policies.PolicyCreate{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ProviderId:  "fc4d5750-22e7-4798-8a46-f48f62c4c1da",
			Parameters: policies.PolicyParam{
				Common: map[string]interface{}{},
			},
			ScheduledOperations: resourceVBSScheduleV2(d),

			Resources: resourceVBSResourceV2(d),
		},
	}
	log.Printf("[INFO] create opts to create: %#v ",createOpts)
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)
	create, err := policies.Create(policyClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup: %s", err)
	}
	log.Printf("[DEBUG] Create: %#v %s", create, create.Id)
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
		return fmt.Errorf("Error waiting for Backup Policy (%s) to become available: %s",	create.Id, StateErr)
	}

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Backup (%s) to become available", create.Id)

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
	for _, hosts := range n.Resources{
		mapping := map[string]interface{}{
			"id":		  hosts.Id,
			"type":         hosts.Type,
			"r_name": hosts.Name,

		}
		resourcelist = append(resourcelist, mapping)
	}

	var scheduledlist []map[string]interface{}
	for _, hosts := range n.ScheduledOperations{
		mapping := map[string]interface{}{
			"so_description":		  hosts.Description,
			"enabled":         hosts.Enabled,
			//"trigger_id":         hosts.TriggerId,
			"so_name":         hosts.Name,
			"operation_type":         hosts.OperationType,
			"max_backups":         hosts.OperationDefinition.MaxBackups,
			//"retention_duration_days":         hosts.OperationDefinition.RetentionDurationDays,
			//"permanent":         hosts.OperationDefinition.Permanent,
			//"plan_id":         hosts.OperationDefinition.PlanId,
			//"provider_id":         hosts.OperationDefinition.ProviderId,
			"scheduler_id":         hosts.Trigger.Id,
			"scheduler_name":         hosts.Trigger.Name,
			"scheduler_type":         hosts.Trigger.Type,
			"pattern":         hosts.Trigger.Properties.Pattern,
			"id": hosts.Id,

		}
		scheduledlist = append(scheduledlist, mapping)
	}



	d.Set("description", n.Description)
	d.Set("id", n.Id)
	d.Set("name", n.Name)
	d.Set("common", n.Parameters.Common)
	d.Set("project_id", n.ProjectId)
	d.Set("provider_id", n.ProviderId)
	//d.Set("id", n.Resources[0].Id)
	//d.Set("type", n.Resources[0].Type)
	//d.Set("name", n.Resources[0].Name)
	//d.Set("extra_info", n.Resources[0].ExtraInfo)
	//d.Set("description", n.ScheduledOperations[0].Description)
	//d.Set("enabled", n.ScheduledOperations[0].Enabled)
	//d.Set("trigger_id", n.ScheduledOperations[0].TriggerId)
	//d.Set("name", n.ScheduledOperations[0].Name)
	//d.Set("operation_type", n.ScheduledOperations[0].OperationType)
	//d.Set("max_backups", n.ScheduledOperations[0].OperationDefinition.MaxBackups)
	//d.Set("retention_duration_days", n.ScheduledOperations[0].OperationDefinition.RetentionDurationDays)
	//d.Set("permanent", n.ScheduledOperations[0].OperationDefinition.Permanent)
	//d.Set("plan_id", n.ScheduledOperations[0].OperationDefinition.PlanId)
	//d.Set("pattern", n.ScheduledOperations[0].Trigger.Properties.Pattern)
	//d.Set("type", n.ScheduledOperations[0].Trigger.Type)
	//d.Set("id", n.ScheduledOperations[0].Trigger.Id)
	//d.Set("name", n.ScheduledOperations[0].Trigger.Name)
	//d.Set("id", n.ScheduledOperations[0].Id)
	d.Set("status", n.Status)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("resources", resourcelist); err != nil {
		return err
	}
	if err := d.Set("scheduled_operations", scheduledlist); err != nil {
		return err
	}



	return nil
}

func resourceCSBSBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Share File: %s", err)
	}
	var updateOpts policies.UpdateOpts

	updateOpts.Policy.Name = d.Get("name").(string)
	updateOpts.Policy.Description = d.Get("description").(string)
	//updateOpts.Policy.Resources = resourceVBSResourceUpdateV2(d)
	updateOpts.Policy.ScheduledOperations = resourceVBSScheduleUpdateV2(d)

	log.Printf("[DEBUG] resourceVBSScheduleUpdateV2 :  %#v", resourceVBSScheduleUpdateV2(d))
	_, err = policies.Update(policyClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud RTS Software Deployment: %s", err)
	}
	log.Printf("[DEBUG] err :  %#v",  err)
	return resourceCSBSBackupPolicyRead(d, meta)
}



func resourceCSBSBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.backupV1Client(GetRegion(d, config))
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

func waitForCSBSBackupPolicyActive(policyClient *golangsdk.ServiceClient, shareID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := policies.Get(policyClient, shareID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "error" {
			return n, n.Status, nil
		}
		return n, n.Status, nil
	}
}

func waitForVBSPolicyDelete(vbsClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := policies.Get(vbsClient, policyID).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", policyID)
				return r, "deleted", nil
			}
			return r, "available", err
		}

		policy := policies.Delete(vbsClient, policyID)
		err = policy.Err
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





func resourceVBSScheduleV2(d *schema.ResourceData) []policies.ScheduledOperations {
	rawTags := d.Get("scheduled_operations").([]interface{})
	tags := make([]policies.ScheduledOperations, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.ScheduledOperations{
			Name:   rawMap["so_name"].(string),
			Description: rawMap["so_description"].(string),
			Enabled: rawMap["enabled"].(bool),
			OperationType:rawMap["operation_type"].(string),
			Trigger: policies.Trigger{
				Properties : policies.TriggerProperties{
					Pattern : rawMap["pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinition{
				MaxBackups: rawMap["max_backups"].(string),
			},

		}
	}
	return tags
}

func resourceVBSResourceV2(d *schema.ResourceData) []policies.Resource {
	rawTags := d.Get("resources").([]interface{})
	tags := make([]policies.Resource, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.Resource{
			Name:   rawMap["r_name"].(string),
			Id: rawMap["id"].(string),
			Type: rawMap["type"].(string),

		}
	}
	return tags
}




func resourceVBSScheduleUpdateV2(d *schema.ResourceData) []policies.ScheduledOperationsUpdate {
	rawTags := d.Get("scheduled_operations").([]interface{})
	tags := make([]policies.ScheduledOperationsUpdate, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.ScheduledOperationsUpdate{
			Id:		rawMap["id"].(string),
			Name:   rawMap["so_name"].(string),
			Description: rawMap["so_description"].(string),
			Enabled: rawMap["enabled"].(bool),
			//TriggerId: rawMap["trigger_id"].(string),
			Trigger: policies.TriggerUpdate{
				Properties : policies.TriggerPropertiesUpdate{
					Pattern : rawMap["pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinitionUpdate{
				MaxBackups: rawMap["max_backups"].(string),
			},

		}
	}
	return tags
}

func resourceVBSResourceUpdateV2(d *schema.ResourceData) []policies.ResourceUpdate {
	rawTags := d.Get("resources").([]interface{})
	tags := make([]policies.ResourceUpdate, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = policies.ResourceUpdate{
			Name:   rawMap["r_name"].(string),
			Id: rawMap["id"].(string),
			Type: rawMap["type"].(string),

		}
	}
	return tags
}