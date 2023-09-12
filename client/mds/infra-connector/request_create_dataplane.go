package infra_connector

type DataPlaneCreateRequest struct {
	AccountId     string `json:"accountId"`
	CertificateId string `json:"certificateId"`
	Name          string `json:"name"`
	TshirtSize    string `json:"nodePoolType"`
	Region        string `json:"region"`
}
