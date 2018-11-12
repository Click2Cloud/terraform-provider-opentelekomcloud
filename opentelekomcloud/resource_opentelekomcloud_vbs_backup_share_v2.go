package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/shares"

	"github.com/huaweicloud/golangsdk"
)

func resourceVBSBackupShareV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVBSBackupShareV2Create,
		Read:   resourceVBSBackupShareV2Read,
		Delete: resourceVBSBackupShareV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"backup_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"to_project_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"container": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_metadata": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceBackupShareToProjectIdsV2(d *schema.ResourceData) []string {
	rawProjectIDs := d.Get("to_project_ids").(*schema.Set)
	projectids := make([]string, rawProjectIDs.Len())
	for i, raw := range rawProjectIDs.List() {
		projectids[i] = raw.(string)
	}
	return projectids
}

func resourceVBSBackupShareV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vbs client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ToProjectIDs: resourceBackupShareToProjectIdsV2(d),
		BackupID:     d.Get("backup_id").(string),
	}

	n, err := shares.Create(vbsClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud VBS Backup Share: %s", err)
	}

	share := n[0]
	d.SetId(share.BackupID)

	log.Printf("[INFO] VBS Backup Share ID: %s", d.Id())

	return resourceVBSBackupShareV2Read(d, meta)
}

func resourceVBSBackupShareV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Vbs client: %s", err)
	}

	backups, err := shares.List(vbsClient, shares.ListOpts{BackupID: d.Id()})
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Vbs: %s", err)
	}

	n := backups[0]

	d.Set("backup_id", n.BackupID)
	d.Set("backup_name", n.Backup.Name)
	d.Set("backup_status", n.Backup.Status)
	d.Set("description", n.Backup.Description)
	d.Set("availability_zone", n.Backup.AvailabilityZone)
	d.Set("volume_id", n.Backup.VolumeID)
	d.Set("size", n.Backup.Size)
	d.Set("service_metadata", n.Backup.ServiceMetadata)
	d.Set("container", n.Backup.Container)
	d.Set("snapshot_id", n.Backup.SnapshotID)
	d.Set("region", GetRegion(d, config))
	d.Set("to_project_ids", resourceToProjectIdsV2(backups))
	d.Set("share_ids", resourceShareIDsV2(backups))

	return nil
}

func resourceVBSBackupShareV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.vbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Vbs client: %s", err)
	}

	deleteopts := shares.DeleteOpts{IsBackupID: true}

	err = shares.Delete(vbsClient, d.Id(), deleteopts).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud Vbs Backup Share %s", d.Id())

		}
		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				log.Printf("[INFO] Error deleting OpenTelekomCloud Vbs Backup Share %s", d.Id())
			}
		}
		log.Printf("[INFO] Successfully deleted OpenTelekomCloud Vbs Backup Share %s", d.Id())
	}

	d.SetId("")
	return nil
}

func resourceToProjectIdsV2(s []shares.Share) []string {
	projectids := make([]string, len(s))
	for i, raw := range s {
		projectids[i] = raw.ToProjectID
	}
	return projectids
}

func resourceShareIDsV2(s []shares.Share) []string {
	shareids := make([]string, len(s))
	for i, raw := range s {
		shareids[i] = raw.ID
	}
	return shareids
}
