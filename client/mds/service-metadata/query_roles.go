package service_metadata

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MDSRolesQuery struct {
	Type string `schema:"serviceType"`
	model.PageQuery
}
