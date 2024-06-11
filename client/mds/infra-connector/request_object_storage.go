package infra_connector

type ObjectStorageCreateRequest struct {
	Name            string `json:"name"`
	BucketName      string `json:"bucketName"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type ObjectStorageUpdateRequest struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}
