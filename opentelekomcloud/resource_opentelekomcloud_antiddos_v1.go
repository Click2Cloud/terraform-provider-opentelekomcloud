package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/huaweicloud/golangsdk"
)

func resourceAntiDdosV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceAntiDdosV1Create,
		Read:   resourceAntiDdosV1Read,
		Update: resourceAntiDdosV1Update,
		Delete: resourceAntiDdosV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_l7": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"traffic_pos_id": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAntiDdosTrafficPosID,
			},
			"http_request_pos_id": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAntiDdosHttpRequestPosID,
			},
			"cleaning_access_pos_id": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAntiDdosCleaningAccessPosID,
			},
			"app_type_id": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAntiDdosAppTypeID,
			},
			"floating_ip_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAntiDdosV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating AntiDdos client: %s", err)
	}

	createOpts := antiddos.CreateOpts{
		EnableL7:            d.Get("enable_l7").(bool),
		TrafficPosId:        d.Get("traffic_pos_id").(int),
		HttpRequestPosId:    d.Get("http_request_pos_id").(int),
		CleaningAccessPosId: d.Get("cleaning_access_pos_id").(int),
		AppTypeId:           d.Get("app_type_id").(int),
	}

	_, err = antiddos.Create(antiddosClient, d.Get("floating_ip_id").(string), createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating AntiDdos: %s", err)
	}

	d.SetId(d.Get("floating_ip_id").(string))

	log.Printf("[INFO] AntiDdos ID: %s", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"configging"},
		Target:     []string{"normal"},
		Refresh:    waitForAntiDdosActive(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for AntiDdos (%s) to become normal: %s",
			d.Id(), stateErr)
	}

	return resourceAntiDdosV1Read(d, meta)

}

func resourceAntiDdosV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating AntiDdos client: %s", err)
	}

	n, err := antiddos.Get(antiddosClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving AntiDdos: %s", err)
	}

	d.Set("enable_l7", n.EnableL7)
	d.Set("app_type_id", n.AppTypeId)
	d.Set("cleaning_access_pos_id", n.CleaningAccessPosId)
	d.Set("traffic_pos_id", n.TrafficPosId)
	d.Set("http_request_pos_id", n.HttpRequestPosId)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceAntiDdosV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating AntiDdos client: %s", err)
	}

	var updateOpts antiddos.UpdateOpts

	updateOpts.EnableL7 = d.Get("enable_l7").(bool)
	updateOpts.AppTypeId = d.Get("app_type_id").(int)
	updateOpts.CleaningAccessPosId = d.Get("cleaning_access_pos_id").(int)
	updateOpts.TrafficPosId = d.Get("traffic_pos_id").(int)
	updateOpts.HttpRequestPosId = d.Get("http_request_pos_id").(int)

	_, err = antiddos.Update(antiddosClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating AntiDdos: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"configging"},
		Target:     []string{"normal"},
		Refresh:    waitForAntiDdosActive(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      3 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for AntiDdos (%s) to become normal: %s", stateErr)
	}

	return resourceAntiDdosV1Read(d, meta)
}

/*func resourceAntiDdosV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating AntiDdos client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal","configging"},
		Target:     []string{"notConfig"},
		Refresh:    waitForAntiDdosDelete(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting AntiDdos: %s", err)
	}

	d.SetId("")
	return nil
}*/

func resourceAntiDdosV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}

	_, err = antiddos.Delete(antiddosClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud tags: %s", err)
	}

	time.Sleep(60 * time.Second)
	d.SetId("")
	return nil
}

func waitForAntiDdosActive(antiddosClient *golangsdk.ServiceClient, antiddosId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := antiddos.Get(antiddosClient, antiddosId).Extract()
		if err != nil {
			return nil, "", err
		}

		return s, "normal", nil
	}
}

func waitForAntiDdosDelete(antiddosClient *golangsdk.ServiceClient, antiddosId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := antiddos.Get(antiddosClient, antiddosId).Extract()
		log.Print("[DEBUG] Get antiddos %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted AntiDdos %s", antiddosId)
				return r, "notConfig", nil
			}
			return r, "normal", err
		}

		n, err := antiddos.Delete(antiddosClient, antiddosId).Extract()
		log.Print("[DEBUG] Delete antiddos %#v", n)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] error 404")
				log.Printf("[INFO] Successfully deleted AntiDdos %s", antiddosId)
				return r, "notConfig", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					log.Printf("[DEBUG] error 409")
					return r, "normal", nil
				}
			}
			return r, "normal", err
		}

		return r, "normal", nil
	}
}

/*func waitForAntiDdosDelete(antiddosClient *golangsdk.ServiceClient, antiddosId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		//r, err := antiddos.Get(antiddosClient, antiddosId).Extract()
		ddosstatus, err := antiddos.ListStatus(antiddosClient, antiddos.ListStatusOpts{FloatingIpId: antiddosId})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted AntiDdos %s", antiddosId)
				return ddosstatus, "notConfig", nil
			}
			return ddosstatus, "normal", err
		}
		r := ddosstatus[0]
		if r.Status != "configging" {
			err := backups.Delete(antiddosClient, antiddosId).ExtractErr()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[INFO] Successfully deleted VBS backup %s", antiddosId)
					return r, "notConfig", nil
				}
				return r, r.Status, err
			}
		}

		return r, "normal", nil
	}
}*/
