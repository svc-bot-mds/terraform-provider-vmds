package model

type MdsClusterMetaData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	ServiceType string `json:"serviceType"`
	Status      string `json:"status"`
}
