package infra_connector

type CloudAccountCreateRequest struct {
	ProviderType string          `json:"type"`
	Name         string          `json:"name"`
	Shared       bool            `json:"shared"`
	Credentials  CredentialModel `json:"credentials"`
	Tags         []string        `json:"tags"`
}

type CloudAccountUpdateRequest struct {
	Credentials CredentialModel `json:"credentials"`
	Tags        []string        `json:"tags"`
}

type CredentialModel struct {
	// gcp credentials
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

	// aws credentials
	ACCESSKEYID                      string `json:"ACCESS_KEY_ID,omitempty"`
	SECRETACCESSKEY                  string `json:"SECRET_ACCESS_KEY,omitempty"`
	TargetAmazonAccountId            string `json:"targetAmazonAccountId,omitempty"`
	PowerUserRoleNameInTargetAccount string `json:"powerUserRoleNameInTargetAccount,omitempty"`

	// vrli credentials
	ApiEndpoint string `json:"apiEndpoint,omitempty"`
	ApiKey      string `json:"apiKey,omitempty"`

	// TKGS credentials
	Username         string `json:"userName,omitempty"`
	Password         string `json:"password,omitempty"`
	SupervisorMgmtIP string `json:"supervisorManagementIP,omitempty"`
	VsphereNamespace string `json:"vsphereNamespace,omitempty"`

	// TKGM credentials
	KubeConfigBase64 string `json:"kubeconfigBase64,omitempty"`

	// Openshift credentials
	//Username         string `json:"userName,omitempty"`
	//Password         string `json:"password,omitempty"`
	Domain string `json:"domain,omitempty"`

	// TAS credentials
	//Username         string `json:"userName,omitempty"`
	//Password         string `json:"password,omitempty"`
	OperationManagerIP string `json:"operationManagerIp,omitempty"`
	CfUsername         string `json:"cfUserName,omitempty"`
	CfPassword         string `json:"cfPassword,omitempty"`
	CfApiHost          string `json:"cfApiHost,omitempty"`
}
