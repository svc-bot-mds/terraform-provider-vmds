package customer_metadata

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsPoliciesQuery struct {
	Id         string   `schema:"id,omitempty"`
	Type       string   `schema:"serviceType"`
	Names      []string `schema:"name,omitempty"`
	ResourceId string   `schema:"resourceId,omitempty"`
	Name       string   `schema:"name,omitempty"`
	model.PageQuery
}
