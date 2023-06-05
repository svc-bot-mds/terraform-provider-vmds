package model

type MdsServiceAccount struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Status string   `json:"status,omitempty"`
	Tags   []string `json:"tags"`
}

type MdsServiceAccountCreate struct {
	OAuthCredentials []*MdsServiceAccountAuthCredentials `json:"oauthCredentials,omitempty"`
}
type MdsServiceAccountAuthCredentials struct {
	UserName   string                        `json:"username,omitempty"`
	Credential *MdsServiceAccountCredentials `json:"credential,omitempty"`
}
type MdsServiceAccountCredentials struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	GrantType    string `json:"grantType,omitempty"`
	OrgId        string `json:"orgId,omitempty"`
}

type MDSServieAccountOauthApp struct {
	AppId       string      `json:"appId"`
	AppType     string      `json:"appType"`
	Created     string      `json:"created"`
	CreatedBy   string      `json:"createdBy"`
	Description string      `json:"description"`
	Modified    string      `json:"modified"`
	ModifiedBy  string      `json:"modifiedBy"`
	TTLSpec     *MDSTTLSpec `json:"ttlSpec"`
}

type MDSTTLSpec struct {
	Description string `json:"description,omitempty"`
	TimeUnit    string `json:"timeUnit"`
	TTL         int64  `json:"ttl"`
}
