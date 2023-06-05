package customer_metadata

type MdsCreateSvcAccountRequest struct {
	AccountType string   `json:"accountType"`
	Usernames   []string `json:"usernames"`
	PolicyIds   []string `json:"policyIds"`
	Tags        []string `json:"tags"`
}

type MDSOauthAppUpdateRequest struct {
	Description string `json:"description,omitempty"`
	TimeUnit    string `json:"timeUnit"`
	TTL         int64  `json:"ttl"`
}
