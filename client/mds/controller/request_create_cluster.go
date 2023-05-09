package controller

type MdsClusterCreateRequest struct {
	Name             string   `json:"name"`
	ServiceType      string   `json:"serviceType"`
	Provider         string   `json:"provider"`
	InstanceSize     string   `json:"instanceSize"`
	Region           string   `json:"region"`
	Dedicated        bool     `json:"dedicated"`
	Tags             []string `json:"tags,omitempty"`
	NetworkPolicyIds []string `json:"networkPolicyIds,omitempty"`
}
