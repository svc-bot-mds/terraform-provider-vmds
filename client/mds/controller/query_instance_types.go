package controller

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsInstanceTypesQuery struct {
	ServiceType string `schema:"serviceType"`
	model.PageQuery
}
