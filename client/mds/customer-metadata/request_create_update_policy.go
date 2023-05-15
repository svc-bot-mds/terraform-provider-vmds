package customer_metadata

type MdsCreateUpdatePolicyRequest struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	ServiceType     string              `json:"serviceType"`
	PermissionsSpec []MdsPermissionSpec `json:"permissionsSpec,omitempty"`
	NetworkSpecs    []MdsNetworkSpecs   `json:"networkSpecs,omitempty"`
}

type MdsPermissionSpec struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
}

type MdsNetworkSpecs struct {
	Cidr           string   `json:"cidr"`
	NetworkPortIds []string `json:"networkPortIds"`
}
