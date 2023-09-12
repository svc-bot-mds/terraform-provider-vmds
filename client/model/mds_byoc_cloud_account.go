package model

type MdsCloudAccount struct {
	Id          string `json:"id"`
	Email       string `json:"userEmail"`
	Name        string `json:"name"`
	AccountType string `json:"accountType"`
	OrgId       string `json:"orgId"`
}
