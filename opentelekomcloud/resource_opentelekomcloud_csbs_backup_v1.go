package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
)

func resourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSBSBackupV1Create,
		Read:   resourceCSBSBackupV1Read,
		Delete: resourceCSBSBackupV1Delete,
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
			"backup_record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OS::Nova::Server",
				ForceNew: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
			"extra_info": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_graph": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_plan": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"extra_info": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceCSBSBackupV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	backupClient, err := config.backupV1Client(GetRegion(d, config))

	log.Printf("[DEBUG] queryOpts: %s", backupClient)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup Client: %s", err)
	}
	queryOpts := backup.ResourceBackupCapOpts{
		CheckProtectable: []backup.ResourceCapQueryParams{
			{
				ResourceId:   d.Get("resource_id").(string),
				ResourceType: d.Get("resource_type").(string),
			},
		},
	}

	log.Printf("[DEBUG] queryOpts: %s", queryOpts)

	query, err := backup.QueryResourceBackupCapability(backupClient, queryOpts).ExtractQueryResponse()
	log.Printf("[DEBUG] query backup: %s", query[0].ResourceId)
	if query[0].Result == true {

		createOpts := backup.CreateOpts{
				BackupName:   d.Get("backup_name").(string),
				Description:  d.Get("description").(string),
				ResourceType: d.Get("resource_type").(string),
				//ExtraInfo:    d.Get("extra_info").(string),
				Tags:         resourceCSBSTagsV1(d),

		}
		log.Printf("[DEBUG] createOpts: %s", createOpts)
		create, err := backup.Create(backupClient, query[0].ResourceId, createOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud backup: %s", err)
		}
		log.Printf("[DEBUG] create: %#v", create)
		backupOpts := backup.ListOpts{CheckpointId: create.Id}
		backupItems, err := backup.List(backupClient, backupOpts)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error retrieving OpenTelekomCloud Backup : %s", err)
		}

		n := backupItems[0]

		d.SetId(n.Id)
		d.Set("backup_record_id", create.Id)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"protecting"},
			Target:     []string{"available"},
			Refresh:    waitForCSBSBackupActive(backupClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Minute,
			MinTimeout: 3 * time.Minute,
		}
		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for Backup (%s) to become Available: %s",
				create.Id, stateErr)
		}

		log.Printf("[DEBUG] Waiting for OpenTelekomCloud Backup (%s) to become available", create.Id)
	} else {
		return fmt.Errorf("Server (%s) is already in service : %s",
			query[0].ResourceId, query[0].ErrorMsg)
	}

	return resourceCSBSBackupV1Read(d, meta)

}

func resourceCSBSBackupV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	backupClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup: %s", err)
	}

	n, err := backup.Get(backupClient, d.Id()).ExtractBackup()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud backup: %s", err)
	}

	d.Set("id", n.CheckpointId)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCSBSBackupV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	backupClient, err := config.backupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Backup: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCSBSBackupDelete(backupClient, d.Id(), d.Get("backup_record_id").(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud backup: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCSBSBackupActive(backupClient *golangsdk.ServiceClient, checkpointItemID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := backup.Get(backupClient, checkpointItemID).ExtractBackup()
		if err != nil {
			return nil, "", err
		}

		if n.Id == "error" {
			return n, n.Status, nil
		}
		return n, n.Status, nil
	}
}

func waitForCSBSBackupDelete(backupClient *golangsdk.ServiceClient, backupId string, backupRecordID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := backup.Get(backupClient, backupId).ExtractBackup()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Backup %s", backupId)
				return r, "deleted", nil
			}
			return r, "deleting", err
		}

		backups := backup.Delete(backupClient, backupRecordID)

		if backups.Err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud Backup %s", backupId)
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "deleting", nil
				}
			}
			return r, "deleting", err
		}

		return r, r.Status, nil
	}
}

func resourceCSBSTagsV1(d *schema.ResourceData) []backup.ResourceTag {
	rawTags := d.Get("tags").([]interface{})
	tags := make([]backup.ResourceTag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tags[i] = backup.ResourceTag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tags
}
