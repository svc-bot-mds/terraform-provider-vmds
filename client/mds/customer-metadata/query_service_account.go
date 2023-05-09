package customer_metadata

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsServiceAccountsQuery struct {
	AccountType string   `schema:"accountType"`
	Name        []string `schema:"name,omitempty"`
	model.PageQuery
}
