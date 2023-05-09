package model

type MdsUser struct {
	Id           string        `json:"id"`
	Email        string        `json:"email"`
	Name         string        `json:"name"`
	Status       string        `json:"status"`
	OrgRoles     []MdsRoleMini `json:"orgRoles,omitempty"`
	ServiceRoles []MdsRoleMini `json:"serviceRoles"`
	Tags         []string      `json:"tags"`
}
