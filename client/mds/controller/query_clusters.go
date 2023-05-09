package controller

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsClustersQuery struct {
	ServiceType   string `schema:"serviceType"`
	Name          string `schema:"name,omitempty"`
	FullNameMatch bool   `schema:"MATCH_FULL_WORD,omitempty"`
	model.PageQuery
}
