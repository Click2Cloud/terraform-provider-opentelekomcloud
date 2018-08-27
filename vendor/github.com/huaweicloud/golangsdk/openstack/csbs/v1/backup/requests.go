package backup

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// SortDir is a type for specifying in which direction to sort a list of Backup.
type SortDir string

var (
	// SortAsc is used to sort a list of Shares in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of Shares in descending order.
	SortDesc SortDir = "desc"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the attributes you want to see returned. SortKey allows you to
// sort by a particular  attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status       string `q:"status"`
	Limit        string `q:"limit"`
	Marker       string `q:"marker"`
	Sort         string `q:"sort"`
	Name         string `q:"name"`
	ResourceId   string `q:"resource_id"`
	CheckpointId string `q:"checkpoint_id"`
	ID           string
	ResourceType string `q:"resource_type"`
	Description  string
}

// List returns collection of
// backups. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]CheckpointItem, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allBackups, err := ExtractBackups(pages)
	if err != nil {
		return nil, err
	}

	return FilterBackups(allBackups, opts)
}

func FilterBackups(backups []CheckpointItem, opts ListOpts) ([]CheckpointItem, error) {

	var refinedBackups []CheckpointItem
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["Id"] = opts.ID
	}
	if opts.Description != "" {
		m["Description"] = opts.Description
	}

	if len(m) > 0 && len(backups) > 0 {
		for _, backup := range backups {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&backup, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedBackups = append(refinedBackups, backup)
			}
		}

	} else {
		refinedBackups = backups
	}

	return refinedBackups, nil
}

func getStructField(v *CheckpointItem, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a Backup. This object is
// passed to backup.Create().
type CreateOpts struct {
	Protect ProtectParam `json:"protect" required:"true"`
}

type ProtectParam struct {
	BackupName   string        `json:"backup_name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ResourceType string        `json:"resource_type,omitempty"`
	Tags         []ResourceTag `json:"tags,omitempty"`
	ExtraInfo    string        `json:"extra_info,omitempty"`
}

// ToBackupCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new backup based on the values in CreateOpts. To extract
// the Backup object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, resourceid string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, resourceid), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// QueryOptsBuilder allows extensions to add additional parameters to the
// QueryResourceCreate request.
type QueryOptsBuilder interface {
	ToQueryResourceCreateMap() (map[string]interface{}, error)
}

// QueryResourceOpts contains the options for querying whether resources can be backed up. This object is
// passed to backup.QueryResourceCreate().
type QueryResourceOpts struct {
	CheckProtectable []ProtectableParam `json:"check_protectable" required:"true"`
}

type ProtectableParam struct {
	ResourceId   string `json:"resource_id" required:"true"`
	ResourceType string `json:"resource_type" required:"true"`
}

// ToQueryResourceCreateMap assembles a request body based on the contents of a
// QueryResourceOpts.
func (opts QueryResourceOpts) ToQueryResourceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// QueryResourceCreate will query whether resources can be backed up based on the values in QueryResourceOpts. To extract
// the Backup object from the response, call the Extract method on the
// QueryResult.
func QueryResourceCreate(client *golangsdk.ServiceClient, opts QueryOptsBuilder) (r QueryResult) {
	b, err := opts.ToQueryResourceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(resourceURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get will get a single backup with specific ID.
func Get(client *golangsdk.ServiceClient, checkpoint_item_id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, checkpoint_item_id), &r.Body, nil)

	return

}

// Delete will delete an existing backup.
func Delete(client *golangsdk.ServiceClient, checkpoint_id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, checkpoint_id), &golangsdk.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: nil,
	})
	return
}
