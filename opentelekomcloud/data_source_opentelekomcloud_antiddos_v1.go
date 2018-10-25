package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
)

func dataSourceAntiDdosV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntiDdosV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"floating_ip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"floating_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period_start": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bps_attack": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bps_in": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_bps": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"pps_in": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"pps_attack": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_pps": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_eips": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logs": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trigger_bps": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trigger_pps": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trigger_http_pps": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAntiDdosV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))

	listStatusOpts := antiddos.ListStatusOpts{
		FloatingIpId: d.Get("floating_ip_id").(string),
		Status:       d.Get("status").(string),
		Ip:           d.Get("floating_ip_address").(string),
	}

	refinedAntiddos, err := antiddos.ListStatus(antiddosClient, listStatusOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the defense status of  EIP: %s", err)
	}

	if len(refinedAntiddos) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedAntiddos) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	ddosStatus := refinedAntiddos[0]

	log.Printf("[INFO] Retrieved defense status of  EIP %s using given filter", ddosStatus.FloatingIpId)

	d.SetId(ddosStatus.FloatingIpId)

	d.Set("floating_ip_id", ddosStatus.FloatingIpId)
	d.Set("floating_ip_address", ddosStatus.FloatingIpAddress)
	d.Set("network_type", ddosStatus.NetworkType)
	d.Set("status", ddosStatus.Status)

	d.Set("region", GetRegion(d, config))

	traffic, err := antiddos.DailyReport(antiddosClient, ddosStatus.FloatingIpId).Extract()
	log.Printf("traffic %#v", traffic)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the traffic of a specified EIP: %s", err)
	}

	//dailyTraffic := traffic[0]

	period_start := make([]int, 0)
	for _, route := range traffic {
		period_start = append(period_start, route.PeriodStart)
	}
	d.Set("period_start", period_start)

	bps_in := make([]int, 0)
	for _, route := range traffic {
		bps_in = append(bps_in, route.BpsIn)
	}
	d.Set("bps_in", bps_in)

	bps_attack := make([]int, 0)
	for _, route := range traffic {
		bps_attack = append(bps_attack, route.BpsAttack)
	}
	d.Set("bps_attack", bps_attack)

	total_bps := make([]int, 0)
	for _, route := range traffic {
		total_bps = append(total_bps, route.TotalBps)
	}
	d.Set("total_bps", total_bps)

	pps_in := make([]int, 0)
	for _, route := range traffic {
		pps_in = append(pps_in, route.PpsIn)
	}
	d.Set("pps_in", pps_in)

	pps_attack := make([]int, 0)
	for _, route := range traffic {
		pps_attack = append(pps_attack, route.PpsAttack)
	}
	d.Set("pps_attack", pps_attack)

	total_pps := make([]int, 0)
	for _, route := range traffic {
		total_pps = append(total_pps, route.TotalPps)
	}
	d.Set("total_pps", total_pps)

	listEventOpts := antiddos.ListLogsOpts{}
	event, err := antiddos.ListLogs(antiddosClient, ddosStatus.FloatingIpId, listEventOpts).Extract()
	log.Printf("event %#v", event)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the event of a specified EIP: %s", err)
	}

	d.Set("logs", flattenLogs(event))

	return nil
}

func flattenLogs(logObject *antiddos.Logs) []map[string]interface{} {
	var log []map[string]interface{}

	mapping := map[string]interface{}{
		"start_time":       logObject.StartTime,
		"end_time":         logObject.EndTime,
		"status":           logObject.Status,
		"trigger_bps":      logObject.TriggerBps,
		"trigger_pps":      logObject.TriggerPps,
		"trigger_http_pps": logObject.TriggerHttpPps,
	}

	log = append(log, mapping)

	return log

}
