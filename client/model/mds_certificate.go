package model

type MdsCertificate struct {
	Id         string `json:"id"`
	DomainName string `json:"domainName"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	ExpiryTime string `json:"expirationTime"`
}
