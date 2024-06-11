package model

type MdsObjectStorage struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	BucketName      string `json:"bucketName"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	OrgId           string `json:"orgId"`
	CreatedBy       string `json:"createdBy"`
	ModifiedBy      string `json:"modifiedBy"`
	ExpirationTime  string `json:"expirationTime"`
}
