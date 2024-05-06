package model

type AuthType string

// ClientAuth -
type ClientAuth struct {
	ApiToken     string `json:"apiKey"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
	OrgId        string `json:"orgId"`
	OAuthAppType string `json:"oAuthAppType"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}
