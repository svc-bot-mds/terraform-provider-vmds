package customer_metadata

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsUsersQuery struct {
	AccountType string   `schema:"accountType"`
	Emails      []string `schema:"email,omitempty"`
	Names       []string `schema:"name,omitempty"`
	model.PageQuery
}
