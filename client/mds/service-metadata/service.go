package service_metadata

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/role_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

var (
	defaultPage = &model.PageQuery{
		Index: 0,
		Size:  100,
	}
)

const (
	EndPoint = "servicemetadata"
)

type Service struct {
	*core.Service
}

func NewService(hostUrl *string, root *core.Root) *Service {
	return &Service{
		Service: core.NewService(hostUrl, EndPoint, root),
	}
}

func (s *Service) GetNetworkPorts() ([]model.MDSNetworkPorts, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, MdsServices, NetworkPorts)

	var response []model.MDSNetworkPorts

	_, err := s.Api.Get(&reqUrl, nil, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// GetMdsRoles - Return list of Roles for the users
func (s *Service) GetMdsRoles(query *MDSRolesQuery) (model.MdsRoles, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, MdsServices, Roles)
	var response model.MdsRoles

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}


// GetMdsServiceRoles - Return list of Roles for the service
func (s *Service) GetMdsServiceRoles(query *MDSRolesQuery) (model.MdsServiceRoles, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, MdsServices, Roles)
	var response model.MdsServiceRoles
	query.Type = role_type.RABBITMQ
	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetPolicyTypes - Returns the policy types
func (s *Service) GetPolicyTypes() ([]string, error) {
	urlPath := fmt.Sprintf("%s/%s/%s/%s", s.Endpoint, MdsServices, Policies, Types)
	var response []string

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return response, err
	}

	return response, err
}
