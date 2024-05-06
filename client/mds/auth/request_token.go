package auth

type TokenRequest struct {
	ApiKey        string `json:"apiKey"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
	AccessToken   string `json:"accessToken"`
	OrgId         string `json:"orgId"`
	OAuthAppTypes string `json:"oAuthAppTypes"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}
