package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
)

func dataSourceCSBSBackupV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCSBSBackupV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_record_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_trigger": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"average_speed": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"copy_from": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"copy_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_op": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"incremental": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"progress": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_az": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"space_saving_ratio": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"supported_restore_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_lld": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cloudservicetype": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"imagetype": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"eip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCSBSBackupV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	backupClient, err := config.backupV1Client(GetRegion(d, config))

	listOpts := backup.ListOpts{
		ID:     d.Get("backup_id").(string),
		Name:   d.Get("backup_name").(string),
		Status: d.Get("status").(string),
	}

	refinedbackups, err := backup.List(backupClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve backup: %s", err)
	}

	if len(refinedbackups) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedbackups) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	backups := refinedbackups[0]

	var t []map[string]interface{}
	for _, tags := range backups.Tags {
		mapping := map[string]interface{}{
			"key":   tags.Key,
			"value": tags.Value,
		}
		t = append(t, mapping)
	}

	log.Printf("[INFO] Retrieved Shares using given filter %s: %+v", backups.Id, backups)
	d.SetId(backups.Id)

	d.Set("backup_record_id", backups.CheckpointId)
	d.Set("backup_id", backups.Id)
	d.Set("backup_name", backups.Name)
	d.Set("resource_id", backups.ResourceId)
	d.Set("status", backups.Status)
	d.Set("description", backups.Description)
	d.Set("resource_type", backups.ResourceType)
	d.Set("auto_trigger", backups.ExtendInfo.AutoTrigger)
	d.Set("average_speed", backups.ExtendInfo.AverageSpeed)
	d.Set("copy_from", backups.ExtendInfo.CopyFrom)
	d.Set("copy_status", backups.ExtendInfo.CopyStatus)
	d.Set("fail_op", backups.ExtendInfo.FailOp)
	d.Set("fail_reason", backups.ExtendInfo.FailReason)
	d.Set("image_type", backups.ExtendInfo.ImageType)
	d.Set("incremental", backups.ExtendInfo.Incremental)
	d.Set("progress", backups.ExtendInfo.Progress)
	d.Set("resource_az", backups.ExtendInfo.ResourceAz)
	d.Set("resource_name", backups.ExtendInfo.ResourceName)
	d.Set("size", backups.ExtendInfo.Size)
	d.Set("space_saving_ratio", backups.ExtendInfo.SpaceSavingRatio)
	d.Set("supported_restore_mode", backups.ExtendInfo.SupportedRestoreMode)
	d.Set("support_lld", backups.ExtendInfo.Supportlld)
	d.Set("cloudservicetype", backups.VMMetadata.CloudServiceType)
	d.Set("disk", backups.VMMetadata.Disk)
	d.Set("imagetype", backups.VMMetadata.ImageType)
	d.Set("ram", backups.VMMetadata.Ram)
	d.Set("vcpus", backups.VMMetadata.Vcpus)
	d.Set("eip", backups.VMMetadata.Eip)
	d.Set("private_ip", backups.VMMetadata.PrivateIp)
	if err := d.Set("tags", t); err != nil {
		return err
	}
	d.Set("region", GetRegion(d, config))

	return nil
}
