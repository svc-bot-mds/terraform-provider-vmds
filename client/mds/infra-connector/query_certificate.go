package infra_connector

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

type MDSCertificateQuery struct {
	Name string `json:"name,omitempty"`
	model.PageQuery
}
