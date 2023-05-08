package model

type MdsRoles struct {
	Embedded struct {
		ServiceRoleDTO []struct {
			Roles []struct {
				MdsRoleMini
				Description string `json:"description"`
			} `json:"roles"`
		} `json:"mdsServiceRoleDTOes"`
	} `json:"_embedded"`
}

type MdsRoleMini struct {
	RoleID string `json:"roleId"`
	Name   string `json:"name"`
}
