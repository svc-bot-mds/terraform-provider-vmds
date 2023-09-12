package model

type ByocDataPlane struct {
	Id                   string      `json:"id"`
	Provider             string      `json:"provider"`
	Name                 string      `json:"name"`
	Region               string      `json:"region"`
	K8SVersion           string      `json:"version"`
	Certificate          Certificate `json:"certificate"`
	DataPlaneReleaseName string      `json:"dataPlaneReleaseName"`
	Status               string      `json:"status"`
	TshirtSize           string      `json:"nodePoolType"`
}

type Certificate struct {
	DomainName string `json:"domainName"`
	Name       string `json:"name"`
}
