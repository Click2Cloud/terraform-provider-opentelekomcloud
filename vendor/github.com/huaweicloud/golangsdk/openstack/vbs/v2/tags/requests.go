package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTagsCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create or add a new tag.
type CreateOpts struct {
	Tag ActionTags `json:"tag" required:"true"`
}

// ActionTags is a structure of key value pair, used for create ,
// update and delete operations.
type ActionTags struct {
	//tag key
	Key string `json:"key" required:"true"`
	//tag value
	Value string `json:"value" required:"true"`
}

// ToTagsCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToTagsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//adds a tag to a backup policy based on the values in CreateOpts. To extract
// the tag object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTagsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(commonURL(client, policyID), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

//deletes a tag ,specified in the key for a backup policy.
func Delete(c *golangsdk.ServiceClient, policyID string, key string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, policyID, key), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Get retrieves the tags of a specific backup policy.
func Get(client *golangsdk.ServiceClient, policyID string) (r GetResult) {
	_, r.Err = client.Get(commonURL(client, policyID), &r.Body, nil)
	return
}

// QueryOptsBuilder allows extensions to add additional parameters to the
// Query request.
type QueryOptsBuilder interface {
	ToTagsQueryMap() (map[string]interface{}, error)
}

// QueryOpts contains all the values needed to Query the policy details based on tags.
type QueryOpts struct {
	//Number of queried records. This parameter is not displayed if action is set to count.
	Limit string `json:"limit,omitempty"`
	//Query index,this parameter is not displayed if action is set to count.
	Offset string `json:"offset,omitempty"`
	//List of tags. Backup policies with these tags will be filtered.
	// This list can have a maximum of 10 tags.
	Tags []QueryTags `json:"tags,omitempty"`
	//Backup policies with any tags in this list will be filtered.
	AnyTags []QueryTags `json:"tags_any,omitempty"`
	//List of excluded tags. Backup policies without these tags will be filtered.
	NotTags []QueryTags `json:"not_tags,omitempty"`
	//List of excluded tags. Backups without any of these tags will be filtered.
	NotAnyTags []QueryTags `json:"not_tags_any,omitempty"`
	//Operator. Possible values are filter and count.
	Action string `json:"action" required:"true"`
}
type QueryTags struct {
	////tag key
	Key string `json:"key" required:"true"`
	//list of values
	Values []string `json:"values" required:"true"`
}

// ToTagsQueryMap builds a Query request body from QueryOpts.
func (opts QueryOpts) ToTagsQueryMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Query retrives a backup policy details using tags.To extract
// the tag object from the response, call the ExtractResources method on the
// QueryResults.
func Query(client *golangsdk.ServiceClient, opts QueryOptsBuilder) (r QueryResults) {
	b, err := opts.ToTagsQueryMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(queryURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// BatchOptsBuilder allows extensions to add additional parameters to the
// BatchAction request.
type BatchOptsBuilder interface {
	ToTagsBatchMap() (map[string]interface{}, error)
}

// BatchOpts contains all the values needed to perform BatchAction on the policy tags.
type BatchOpts struct {
	//List of tags to perform batch operation
	Tags []ActionTags `json:"tags,omitempty"`
	//Operator , Possible values are:create, update,delete
	Action ActionType `json:"action" required:"true"`
}

//ActionType specifies the type of batch operation action to be performed
type ActionType string

var (
	// ActionCreate is used to set action operator to create
	ActionCreate ActionType = "create"
	// ActionDelete is used to set action operator to delete
	ActionDelete ActionType = "delete"
	// ActionUpdate is used to set action operator to update
	ActionUpdate ActionType = "update"
)

// ToTagsBatchMap builds a BatchAction request body from BatchOpts.
func (opts BatchOpts) ToTagsBatchMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//BatchAction is used to create ,update or delete the tags of a specified backup policy.
func BatchAction(client *golangsdk.ServiceClient, policyID string, opts BatchOptsBuilder) (r ActionResults) {
	b, err := opts.ToTagsBatchMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, policyID), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
