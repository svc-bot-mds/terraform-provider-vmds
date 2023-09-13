package infra_connector

type CloudAccountCreateRequest struct {
	ProviderType string          `json:"type"`
	Name         string          `json:"name"`
	Credentials  CredentialModel `json:"credentials"`
}

type CredentialModel struct {
	//gcp credentials
	Type                    string `json:"type,omitempty"`
	ProjectId               string `json:"project_id,omitempty"`
	PrivateKeyId            string `json:"private_key_id,omitempty"`
	PrivateKey              string `json:"private_key,omitempty"`
	ClientEmail             string `json:"client_email,omitempty"`
	ClientId                string `json:"client_id,omitempty"`
	AuthUri                 string `json:"auth_uri,omitempty"`
	TokenUri                string `json:"token_uri,omitempty"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertUrl       string `json:"client_x509_cert_url,omitempty"`

	//aws credentials
	ACCESSKEYID                      string `json:"ACCESS_KEY_ID,omitempty"`
	SECRETACCESSKEY                  string `json:"SECRET_ACCESS_KEY,omitempty"`
	TargetAmazonAccountId            string `json:"targetAmazonAccountId,omitempty"`
	PowerUserRoleNameInTargetAccount string `json:"powerUserRoleNameInTargetAccount,omitempty"`

	//vrli credentials
	ApiEndpoint string `json:"apiEndpoint,omitempty"`
	ApiKey      string `json:"apiKey,omitempty"`
}
