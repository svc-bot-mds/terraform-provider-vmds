package infra_connector

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MdsCloudAccountsQuery struct {
	AccountType string `schema:"accountType,omitempty"`
	model.PageQuery
}
