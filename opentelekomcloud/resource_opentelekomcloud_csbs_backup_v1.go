package opentelekomcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
	"fmt"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
	"github.com/huaweicloud/golangsdk"
	"github.com/hashicorp/terraform/helper/resource"
	"log"
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
			"backup_id": {
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
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
			"resource_type":&schema.Schema{
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

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud backup Client: %s", err)
	}
	queryOpts := backup.QueryResourceOpts{
		CheckProtectable: []backup.ProtectableParam{
			{
				ResourceId:	d.Get("resource_id").(string),
				ResourceType:	d.Get("resource_type").(string),
			},
		},
	}

	log.Printf("[DEBUG] queryOpts: %s", queryOpts)

	query, err := backup.QueryResourceCreate(backupClient, queryOpts).ExtractQueryResponse()
	log.Printf("[DEBUG] query backup: %s", query.Protectable[0].ResourceId)


	createOpts := backup.CreateOpts{
		Protect:       backup.ProtectParam{
			BackupName: d.Get("backup_name").(string),
			Description: d.Get("description").(string),
			ResourceType: d.Get("resource_type").(string),
			ExtraInfo: d.Get("extra_info").(string),
		},
	}

	log.Printf("[DEBUG] CreateOpts: %s", createOpts)
	if query.Protectable[0].Result == true {

		create, err := backup.Create(backupClient, query.Protectable[0].ResourceId, createOpts).Extract()

		log.Printf("[DEBUG] create backup: %s", create)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud backup: %s", err)
		}

		backupOpts := backup.ListOpts{CheckpointId: create.Checkpoint.Id}
		backupItems,err := backup.List(backupClient,backupOpts)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error retrieving OpenTelekomCloud Backup Policy: %s", err)
		}

		n := backupItems[0]

		d.SetId(n.Id)
		d.Set("backup_id",create.Checkpoint.Id)
		log.Printf("[DEBUG] set ID: %s", d.Id())

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
				create.Checkpoint.Id, stateErr)
		}

		log.Printf("[DEBUG] Waiting for OpenTelekomCloud Backup (%s) to become available", create.Checkpoint.Id)
	}else {
		return fmt.Errorf("Server (%s) is already in service : %s",
					query.Protectable[0].ResourceId,query.Protectable[0].ErrorMsg)
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
	log.Printf("[DEBUG] get backup: %s", n)

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
		Pending:    []string{"available","deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForCSBSBackupDelete(backupClient, d.Id(),d.Get("backup_id").(string)),
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

func waitForCSBSBackupDelete(backupClient *golangsdk.ServiceClient, backupId string, backupRecordID string ) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := backup.Get(backupClient, backupId).ExtractBackup()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", backupId)
				return r, "deleted", nil
			}
			return r, "deleting", err
		}

		backups := backup.Delete(backupClient, backupRecordID)

		if backups.Err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", backupId)
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

