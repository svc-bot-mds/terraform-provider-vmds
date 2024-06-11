package model

type MdsCloudAccount struct {
	Id             string   `json:"id"`
	Email          string   `json:"userEmail"`
	Name           string   `json:"name"`
	AccountType    string   `json:"accountType"`
	OrgId          string   `json:"orgId"`
	Shared         bool     `json:"shared"`
	Tags           []string `json:"tags"`
	DataPlaneCount int64    `json:"dataplanesCount"`
	Created        string   `json:"created"`
	CreatedBy      string   `json:"createdBy"`
	Modified       string   `json:"modified"`
	ModifiedBy     string   `json:"modifiedBy"`
	ManagementIp   string   `json:"managementIp"`
}
