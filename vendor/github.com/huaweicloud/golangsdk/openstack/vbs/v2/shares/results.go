package shares

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Share struct {
	//Details about the source backup
	Backup Backup `json:"backup"`
	//Backup ID
	BackupId string `json:"backup_id"`
	//Backup share ID
	Id string `json:"id"`
	//ID of the project with which the backup is shared
	ToProjectId string `json:"to_project_id"`
	//ID of the project that shares the backup
	FromProjectId string `json:"from_project_id"`
	// Creation time of the backup share
	CreatedAt string `json:"created_at"`
	//Update time of the backup share
	UpdatedAt string `json:"updated_at"`
	//Whether the backup has been deleted
	Deleted string `json:"deleted"`
	//Deletion time
	DeletedAt string `json:"deleted_at"`
}

type Backup struct {
	//Backup ID
	Id string `json:"id"`
	//Backup name
	Name string `json:"name"`
	//Backup status
	Status string `json:"status"`
	//Backup description
	Description string `json:"description"`
	//AZ where the backup resides
	AZ string `json:"availability_zone"`
	//Source volume ID of the backup
	VolumeId string `json:"volume_id"`
	//Cause of the backup failure
	FailReason string `json:"fail_reason"`
	//Backup size
	Size int `json:"size"`
	//Number of objects on OBS for the disk data
	ObjectCount int `json:"object_count"`
	//Container of the backup
	Container string `json:"container"`
	//Backup creation time
	CreatedAt string `json:"created_at"`
	//Backup metadata
	ServiceMetadata string `json:"service_metadata"`
	//Time when the backup was updated
	UpdatedAt string `json:"updated_at"`
	//Current time
	DataTimeStamp string `json:"data_timestamp"`
	//Whether a dependent backup exists
	DependentBackups bool `json:"has_dependent_backups"`
	//ID of the snapshot associated with the backup
	SnapshotId string `json:"snapshot_id"`
	//Whether the backup is an incremental backup
	Incremental bool `json:"is_incremental"`
}

type commonResult struct {
	golangsdk.Result
}

// SharePage is the page returned by a pager when traversing over a
// collection of Shares.
type SharePage struct {
	pagination.LinkedPageBase
}

// Extract is a function that accepts a result and extracts a share.
func (r commonResult) Extract() ([]Share, error) {
	var s struct {
		Share []Share `json:"shared"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// ExtractGet is a function that accepts a result and extracts a share.
func (r commonResult) ExtractGet() (*Share, error) {
	var s struct {
		Share *Share `json:"shared"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// NextPageURL is invoked when a paginated collection of Shares has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SharePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"shared_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SharePage struct is empty.
func (r SharePage) IsEmpty() (bool, error) {
	is, err := ExtractShares(r)
	return len(is) == 0, err
}

// ExtractShares accepts a Page struct, specifically a SharePage struct,
// and extracts the elements into a slice of Shares struct. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractShares(r pagination.Page) ([]Share, error) {
	var s struct {
		Shares []Share `json:"shared"`
	}
	err := (r.(SharePage)).ExtractInto(&s)
	return s.Shares, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Share.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its ExtractGet
// method to interpret it as a Share.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
