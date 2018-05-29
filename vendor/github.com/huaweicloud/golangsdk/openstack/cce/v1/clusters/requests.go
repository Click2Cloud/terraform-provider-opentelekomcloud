package clusters

import (
	"github.com/huaweicloud/golangsdk"
	"reflect"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// ListOpts allows the filtering  of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned.
type ListOpts struct {
	Name    string `json:"name"`
	ID      string `json:"uuid"`
	AZ      string `json:"az"`
	Type    string `json:"clustertype"`
	VpcName string `json:"vpc"`
	VpcId   string `json:"vpcid"`
	Status  string `json:"status"`
}

// List returns collection of
// clusters. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those clusters that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(client *golangsdk.ServiceClient) (r ListResult) {

	_, r.Err = client.Get(rootURL(client), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})

	return
}

func FilterClusters(clusters []RetrievedCluster, opts ListOpts) ([]RetrievedCluster, error) {

	var refinedClusters []RetrievedCluster
	var matched bool
	m := map[string]FilterStruct{}

	if opts.Name != "" {
		m["Name"] = FilterStruct{Value: opts.Name, Driller: []string{"Metadata"}}
	}
	if opts.ID != "" {
		m["ID"] = FilterStruct{Value: opts.ID, Driller: []string{"Metadata"}}
	}
	if opts.AZ != "" {
		m["AZ"] = FilterStruct{Value: opts.AZ, Driller: []string{"Clusterspec"}}
	}
	if opts.Type != "" {
		m["ClusterType"] = FilterStruct{Value: opts.Type, Driller: []string{"Clusterspec"}}
	}
	if opts.VpcId != "" {
		m["VpcId"] = FilterStruct{Value: opts.VpcId, Driller: []string{"Clusterspec"}}
	}
	if opts.VpcName != "" {
		m["VpcName"] = FilterStruct{Value: opts.VpcName, Driller: []string{"Clusterspec"}}
	}
	if opts.Status != "" {
		m["Status"] = FilterStruct{Value: opts.Status, Driller: []string{"ClusterStatus"}}
	}

	if len(m) > 0 && len(clusters) > 0 {
		for _, cluster := range clusters {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&cluster, key, value.Driller); !(sVal == value.Value) {
					matched = false
				}
			}
			if matched {
				refinedClusters = append(refinedClusters, cluster)
			}
		}

	} else {
		refinedClusters = clusters
	}

	return refinedClusters, nil
}

type FilterStruct struct {
	Value   string
	Driller []string
}

func GetStructNestedField(v *RetrievedCluster, field string, structDriller []string) string {
	r := reflect.ValueOf(v)
	for _, drillField := range structDriller {
		f := reflect.Indirect(r).FieldByName(drillField).Interface()
		r = reflect.ValueOf(f)
	}
	f1 := reflect.Indirect(r).FieldByName(field)
	return string(f1.String())
}

type CreateMetadataspec struct {
	Name string `json:"name" required:"true"`
}

type CreateSpec struct {
	Description     string `json:"description,omitempty"`
	Vpc             string `json:"vpc" required:"true"`
	Subnet          string `json:"subnet" required:"true"`
	Region          string `json:"region" required:"true"`
	SecurityGroupId string `json:"security_group_id,omitempty"`
	ClusterType     string `json:"clustertype,omitempty"`
}

type UpdateSpec struct {
	Description string `json:"description,omitempty"`
	EIP         string `json:"publicip_id,omitempty"`
}

// CreateOpts contains all the values needed to create a new cluster.
type CreateOpts struct {
	Kind       string             `json:"kind" required:"true"`
	ApiVersion string             `json:"apiVersion" required:"true"`
	Metadata   CreateMetadataspec `json:"metadata" required:"true"`
	Spec       CreateSpec         `json:"spec" required:"true"`
}

type CreateOptsBuilder interface {
	ToVpcCreateMap() (map[string]interface{}, error)
}

// ToVpcCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToVpcCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// cluster.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVpcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, nil, reqOpt)
	return
}

// Get retrieves a particular cluster based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// GetCertificate retrieves a particular cluster certificates based on cluster ID.
func GetCertificate(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(certificateURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVpcUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a cluster.
type UpdateOpts struct {
	Kind       string     `json:"kind" required:"true"`
	ApiVersion string     `json:"apiVersion" required:"true"`
	Spec       UpdateSpec `json:"spec" required:"true"`
}

// ToVpcUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update allows clusters to add EIP and description.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular cluster based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
