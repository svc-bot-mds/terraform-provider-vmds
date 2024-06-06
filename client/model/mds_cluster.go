package model

// MdsCluster -
type MdsCluster struct {
	ID                   string              `json:"id,omitempty"`
	OrgId                string              `json:"orgId"`
	Name                 string              `json:"name"`
	ServiceType          string              `json:"serviceType"`
	Provider             string              `json:"provider"`
	InstanceSize         string              `json:"instanceSize"`
	Region               string              `json:"region"`
	Tags                 []string            `json:"tags"`
	Version              string              `json:"version"`
	Status               string              `json:"status"`
	DataPlaneId          string              `json:"dataPlaneId"`
	Metadata             *MdsClusterMetadata `json:"metadata"`
	Created              string              `json:"created"`
	LastUpdated          string              `json:"lastUpdated"`
	IsAuthorized         bool                `json:"isAuthorized,omitempty"`
	MaintenanceStartTime int64               `json:"maintenanceStartTime,omitempty"`
	MaintenanceEndTime   int64               `json:"maintenanceEndTime,omitempty"`
	UpgradeInProgress    bool                `json:"isUpgradeInProgress"`
	PauseUpdates         bool                `json:"pauseUpdates"`
}

type MdsClusterMetadata struct {
	ClusterName      string   `json:"clusterName,omitempty" tfsdk:"cluster_name"`
	ManagerUri       string   `json:"managerUri,omitempty" tfsdk:"manager_uri"`
	ConnectionUri    string   `json:"connectionUri,omitempty" tfsdk:"connection_uri"`
	MetricsEndpoints []string `json:"metricsEnpoints" tfsdk:"metrics_endpoints"`
}
