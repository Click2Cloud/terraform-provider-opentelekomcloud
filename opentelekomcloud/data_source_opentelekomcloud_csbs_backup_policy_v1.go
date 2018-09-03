package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func dataSourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCSBSBackupPolicyV1Read,

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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"resources": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"scheduled_operations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheduled_period_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"scheduled_period_description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"max_backups": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
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
						},
						"pattern": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduled_period_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
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

func dataSourceCSBSBackupPolicyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.backupV1Client(GetRegion(d, config))

	listOpts := policies.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedpolicies, err := policies.List(policyClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve backup policies: %s", err)
	}

	if len(refinedpolicies) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedpolicies) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	policy := refinedpolicies[0]

	var r []map[string]interface{}
	for _, resource := range policy.Resources {
		mapping := map[string]interface{}{
			"resource_id":   resource.Id,
			"resource_type": resource.Type,
			"resource_name": resource.Name,
		}
		r = append(r, mapping)
	}

	var s []map[string]interface{}
	for _, scheduled_operation := range policy.ScheduledOperations {
		mapping := map[string]interface{}{
			"scheduled_period_description": scheduled_operation.Description,
			"enabled":                      scheduled_operation.Enabled,
			"scheduled_period_name":        scheduled_operation.Name,
			"operation_type":               scheduled_operation.OperationType,
			"max_backups":                  scheduled_operation.OperationDefinition.MaxBackups,
			"retention_duration_days":      scheduled_operation.OperationDefinition.RetentionDurationDays,
			"permanent":                    scheduled_operation.OperationDefinition.Permanent,
			"plan_id":                      scheduled_operation.OperationDefinition.PlanId,
			"pattern":                      scheduled_operation.Trigger.Properties.Pattern,
			"scheduler_id":                 scheduled_operation.Trigger.ID,
			"scheduler_name":               scheduled_operation.Trigger.Name,
			"scheduler_type":               scheduled_operation.Trigger.Type,
			"scheduled_period_id":          scheduled_operation.ID,
			"trigger_id":                   scheduled_operation.TriggerID,
		}
		s = append(s, mapping)
	}

	var t []map[string]interface{}
	for _, resource := range policy.Tags {
		mapping := map[string]interface{}{
			"key":   resource.Key,
			"value": resource.Value,
		}
		t = append(t, mapping)
	}

	log.Printf("[INFO] Retrieved Shares using given filter %s: %+v", policy.ID, policy)
	d.SetId(policy.ID)

	d.Set("description", policy.Description)
	d.Set("id", policy.ID)
	d.Set("name", policy.Name)
	d.Set("common", policy.Parameters.Common)
	d.Set("project_id", policy.ProjectId)
	d.Set("provider_id", policy.ProviderId)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("resources", r); err != nil {
		return err
	}
	if err := d.Set("scheduled_operations", s); err != nil {
		return err
	}
	if err := d.Set("tags", t); err != nil {
		return err
	}

	return nil
}
