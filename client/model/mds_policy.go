package model

// MdsPolicy base model for MDS Policy
type MdsPolicy struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	ServiceType     string               `json:"serviceType"`
	ResourceIds     []string             `json:"resourceIds,omitempty"`
	PermissionsSpec []*MdsPermissionSpec `json:"permissionsSpec,omitempty"`
	NetworkSpec     []*MdsNetworkSpec    `json:"networkSpecs,omitempty"`
}
type MdsPermissionSpec struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
}

type MdsNetworkSpec struct {
	CIDR           string   `json:"cidr"`
	NetworkPortIds []string `json:"networkPortIds"`
}
