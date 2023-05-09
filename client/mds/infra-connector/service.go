package infra_connector

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
)

const (
	EndPoint = "infra-connector"
)

type Service struct {
	*core.Service
}

func NewService(hostUrl *string, root *core.Root) *Service {
	return &Service{
		Service: core.NewService(hostUrl, EndPoint, root),
	}
}

func (s *Service) GetRegionsWithDataPlanes(regionsQuery *DataPlaneRegionsQuery) (map[string][]string, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, K8sCluster, Resource)

	var response map[string][]string

	_, err := s.Api.Get(&reqUrl, regionsQuery, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
