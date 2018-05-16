package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/clusters"
	"log"
	"time"
)

func resourceCceClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCceClusterV1Create,
		Read:   resourceCceClusterV1Read,
		Update: resourceCceClusterV1Update,
		Delete: resourceCceClusterV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc:validateClusterName,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"publicip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed:true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed:true,
			},
			"k8s_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional:true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed:true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional:true,
			},
			"vpc_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed:true,
			},
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"az": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"sshkey": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCceClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV1Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of CCE Client: %#v", cceClient)

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud cluster Client: %s", err)
	}

	createOpts := clusters.CreateOpts{
		Kind:       d.Get("kind").(string),
		ApiVersion: d.Get("api_version").(string),
		Metadata:   clusters.CreateMetadataspec{Name: d.Get("name").(string)},
		Spec: clusters.CreateSpec{
			Description:     d.Get("description").(string),
			Vpc:             d.Get("vpc_id").(string),
			Subnet:          d.Get("subnet_id").(string),
			Region:          GetRegion(d, config),
			SecurityGroupId: d.Get("security_group_id").(string),
			ClusterType:     d.Get("type").(string),
		},
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	create := clusters.Create(cceClient, createOpts).ExtractErr()

	if create != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Cluster: %s", create)
	}

	name := d.Get("name").(string)
	listops := clusters.ListOpts{Name: name}

	clusterlist, err := clusters.List(cceClient).ExtractCluster(listops)

	n := clusterlist[0]

	log.Printf("[INFO] cluster ID: %s", n.Metadata.ID)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud CCE cluster (%s) to become available", n.Metadata.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"EMPTY"},
		Refresh:    waitForCceClusterActive(cceClient, n.Metadata.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(n.Metadata.ID)

	if _, ok := d.GetOk("publicip_id"); ok {
		resourceCceClusterV1Update(d, meta)
	}
	return resourceCceClusterV1Read(d, meta)

}

func resourceCceClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Cce client: %s", err)
	}

	n, err := clusters.Get(cceClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Vpc: %s", err)
	}

	log.Printf("[DEBUG] Retrieved cluster %s: %+v", d.Id(), n)

	/*var volumespec []map[string]interface{}
	for _, volume := range n.Clusterspec.ClusterHostList.HostListSpec.HostList[0].Hostspec.Volume {
		mapping := map[string]interface{}{
			"disk_type":        volume.DiskSize,
			"disk_size":        volume.DiskType,
			"volume_type": 		volume.VolumeType,
		}
		volumespec = append(volumespec, mapping)
	}*/

	var hostspec []map[string]interface{}
	for _, hosts := range n.Clusterspec.ClusterHostList.HostListSpec.HostList {
		mapping := map[string]interface{}{
			"name":       hosts.Metadata.Name,
			"id":         hosts.Metadata.ID,
			"private_ip": hosts.Hostspec.PrivateIp,
			"public_ip":  hosts.Hostspec.PublicIp,
			"flavor":     hosts.Hostspec.Flavor,
			"az":         hosts.Hostspec.AZ,
			"sshkey":     hosts.Hostspec.SshKey,
			"status":     hosts.NodeStatus,
			//"volume":     volumespec,
		}
		hostspec = append(hostspec, mapping)
	}

	d.Set("id", n.Metadata.ID)
	d.Set("name", n.Metadata.Name)
	d.Set("status", n.ClusterStatus.Status)
	d.Set("k8s_version", n.K8sVersion)
	d.Set("az", n.Clusterspec.AZ)
	d.Set("cpu", n.Clusterspec.CPU)
	d.Set("type", n.Clusterspec.ClusterType)
	d.Set("vpc_name", n.Clusterspec.VPC)
	d.Set("vpc_id", n.Clusterspec.VpcId)
	d.Set("subnet", n.Clusterspec.Subnet)
	d.Set("endpoint", n.Clusterspec.Endpoint)
	d.Set("external_endpoint", n.Clusterspec.ExternalEndpoint)
	d.Set("security_group_id", n.Clusterspec.SecurityGroupId)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("hosts", hostspec); err != nil {
		return err
	}

	return nil
}

func resourceCceClusterV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE Client: %s", err)
	}

	var updateOpts clusters.UpdateOpts

	updateOpts.Kind = d.Get("kind").(string)
	updateOpts.ApiVersion = d.Get("api_version").(string)

	if d.HasChange("description") {
		updateOpts.Spec.Description = d.Get("description").(string)
	}
	updateOpts.Spec.EIP = d.Get("publicip_id").(string)

	log.Printf("[DEBUG] Updating CCE %s with options: %+v", d.Id(), updateOpts)

	err = clusters.Update(cceClient, d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud CCE: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UPDATING"},
		Target:     []string{"EMPTY","AVAILABLE"},
		Refresh:    waitForCceClusterUpdate(cceClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()


	return resourceVirtualPrivateCloudV1Read(d, meta)
}

func resourceCceClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy CCE cluster: %s", d.Id())

	config := meta.(*Config)
	cceClient, err := config.cceV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING","AVAILABLE","EMPTY"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCceClusterDelete(cceClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud CCE cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceClusterActive(cceClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clusters.Get(cceClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}
		if n.ClusterStatus.Status != "EMPTY"  {
			return n,"CREATING", nil
		}

		return n, n.ClusterStatus.Status, nil
	}
}

func waitForCceClusterUpdate(cceClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clusters.Get(cceClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}
		if n.ClusterStatus.Status != "UPDATING"  {
			return n,n.ClusterStatus.Status, nil
		}

		return n, n.ClusterStatus.Status, nil
	}
}
func waitForCceClusterDelete(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud CCE cluster %s.\n", clusterId)

		r, err := clusters.Get(cceClient, clusterId).Extract()

		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud CCE cluster %s", clusterId)
				return r, "DELETED", nil
			}
			return r, r.ClusterStatus.Status, err
		}

		err = clusters.Delete(cceClient, clusterId).ExtractErr()
		log.Printf("[DEBUG] Value if error: %#v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud CCE cluster %s", clusterId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "AVAILABLE", nil
				}
			}
			return r, "AVAILABLE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud CCE cluster %s still available.\n", clusterId)
		return r, "AVAILABLE", nil
	}
}
