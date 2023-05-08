package customer_metadata

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsPoliciesQuery struct {
	Type       string   `schema:"serviceType"`
	Names      []string `schema:"name,omitempty"`
	ResourceId string   `schema:"resourceId,omitempty"`
	model.PageQuery
}
