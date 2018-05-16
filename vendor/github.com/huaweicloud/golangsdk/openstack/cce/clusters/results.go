package clusters

import (
	"github.com/huaweicloud/golangsdk"

)

type ClusterStatus struct {
	Status    string `json:"status"`
}

type MetaData struct {
	// ID is the unique identifier for the vpc.
	Name 	string `json:"name"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	ID 		string `json:"uuid"`

	//Specifies the range of available subnets in the VPC.
	SpaceID string `json:"spaceuuid"`
}

type ClusterlistSpec struct {
	ClusterHostList 	ClusterHostList `json:"hostList"`
	AZ 					string 		`json:"az"`
	CPU 				int 		`json:"cpu"`
	Memory 				int 		`json:"memory"`
	VPC				 	string		`json:"vpc"`
	VpcId 				string		`json:"vpcid"`
	Subnet 				string		`json:"subnet"`
	Endpoint 			string		`json:"endpoint"`
	ExternalEndpoint 	string		`json:"external_endpoint"`
	SecurityGroupId		string		`json:"security_group_id"`
	ClusterType 		string		`json:"clustertype"`
}

type ClusterHostList struct {
	Kind 				string 		 `json:"kind"`
	ApiVersion 			string 		 `json:"apiVersion"`
	Metadata 			MetaData	 `json:"metadata"`
	HostListSpec  		HostListSpec `json:"spec"`
}

type HostListSpec struct {
	HostList []Host `json:"hostList"`
}

type Host struct {

	Kind string `json:"kind"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	ApiVersion string `json:"apiVersion"`

	//Specifies the range of available subnets in the VPC.
	Metadata MetaData `json:"metadata"`

	Hostspec  HostSpec  `json:"spec"`

	Replicas  int`json:"replicas"`

	NodeStatus  string`json:"status"`

	Message  string`json:"message"`
}

type HostSpec struct {

	ClusterUuid string `json:"clusteruuid"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	ClusterName string `json:"clustername"`

	//Specifies the range of available subnets in the VPC.
	PrivateIp string `json:"privateip"`

	PublicIp  string`json:"publicip"`

	Flavor string `json:"flavor"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	CPU int `json:"cpu"`

	//Specifies the range of available subnets in the VPC.
	Memory int `json:"memory"`

	HostId  string`json:"hostid"`

	AZ string `json:"az"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	Volume []Volume `json:"volume"`

	//Specifies the range of available subnets in the VPC.
	SshKey string `json:"sshkey"`

	HostStatus  HostStatus  `json:"status"`

}

type Volume struct {
	DiskType 	string  `json:"diskType"`
	DiskSize 	int		`json:"diskSize"`
	VolumeType  string	`json:"volumeType"`
}

type HostStatus struct {
	Capacity 		Capacity		`json:"capacity"`
	Allocatable 	Capacity		`json:"allocatable"`
	Conditions 		[]Conditions	`json:"conditions"`
	Addresses 		[]Addresses		`json:"addresses"`
	DaemonEndpoints DaemonEndpoints `json:"daemonEndpoints"`
	NodeInfo 		NodeInfo		`json:"nodeInfo"`
	Images 			[]Images		`json:"images"`
}

type Capacity struct {
	CPU 	string	`json:"cpu"`
	Memory 	string	`json:"memory"`
	Pods 	string	`json:"pods"`
}

type Conditions struct {
	Type 	string	`json:"type"`
	Status	string	`json:"status"`
	Reason  string	`json:"reason"`
}

type Addresses struct {
	Type	string	`json:"type"`
	Address	string	`json:"address"`
}

type DaemonEndpoints struct {
	KubeletEndpoint KubeletEndpoint	`json:"kubeletEndpoint"`
}

type KubeletEndpoint struct {
	Port 	int	`json:"Port"`
}

type NodeInfo struct {
	MachineID				string	`json:"machineID"`
	SystemUUID				string	`json:"systemUUID"`
	BootID					string	`json:"bootID"`
	KernelVersion			string	`json:"kernelVersion"`
	OsImage					string	`json:"osImage"`
	ContainerRuntimeVersion string	`json:"containerRuntimeVersion"`
	KubeletVersion			string	`json:"kubeletVersion"`
	KubeProxyVersion		string	`json:"kubeProxyVersion"`
}

type Images struct {
	Names		[]string	`json:"names"`
	SizeBytes	int			`json:"sizeBytes"`
}

type RetrievedCluster struct {
	// ID is the unique identifier for the vpc.
	Kind string `json:"kind"`
	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	ApiVersion string `json:"apiVersion"`

	//Specifies the range of available subnets in the VPC.
	Metadata MetaData `json:"metadata"`

	// Status indicates whether or not a vpc is currently operational.
	Clusterspec ClusterlistSpec `json:"spec"`

	ClusterStatus ClusterStatus `json:"clusterStatus"`

	K8sVersion  string `json:"k8s_version"`
}

type commonResult struct {
	golangsdk.Result
}


// ListResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type ListResult struct {
	commonResult
}

// ExtractCluster is a function that accepts a result and extracts a cluster.
func (r commonResult) ExtractCluster(opts ListOpts) ([]RetrievedCluster, error) {
	var s []RetrievedCluster
	err := r.ExtractInto(&s)
	if err!= nil{
		return nil,err
	}
	return FilterClusters(s,opts)
}



// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Cluster.
type CreateResult struct {
	golangsdk.ErrResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Cluster.
type UpdateResult struct {
	golangsdk.ErrResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

func (r commonResult) Extract() (*RetrievedCluster, error) {
	var s RetrievedCluster
	err := r.ExtractInto(&s)
	return &s, err
}
