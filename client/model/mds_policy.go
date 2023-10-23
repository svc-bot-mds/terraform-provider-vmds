package model

// MdsPolicy base model for MDS Policy
type MdsPolicy struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	ServiceType     string                `json:"serviceType"`
	ResourceIds     []string              `json:"resourceIds,omitempty"`
	PermissionsSpec []*MdsPermissionsSpec `json:"permissionsSpec,omitempty"`
	NetworkSpec     []*MdsNetworkSpec     `json:"networkSpecs,omitempty"`
}
type MdsPermissionsSpec struct {
	Resource    string            `json:"resource"`
	Permissions []*MdsPermissions `json:"permissions"`
	Role        string            `json:"role"`
}

type MdsNetworkSpec struct {
	CIDR           string   `json:"cidr"`
	NetworkPortIds []string `json:"networkPortIds"`
}

type MdsPermissions struct {
	Name         string `json:"name"`
	PermissionId string `json:"permissionId"`
}
