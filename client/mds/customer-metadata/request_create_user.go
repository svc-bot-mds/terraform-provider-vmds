package customer_metadata

type MdsCreateUserRequest struct {
	AccountType  string         `json:"accountType"`
	Usernames    []string       `json:"usernames"` // List of emails by which to invite/add the users.
	PolicyIds    []string       `json:"policyIds"`
	ServiceRoles []RolesRequest `json:"serviceRoles"`
	Tags         []string       `json:"tags"`
}

type RolesRequest struct {
	RoleId string `json:"roleId,omitempty"`
}
