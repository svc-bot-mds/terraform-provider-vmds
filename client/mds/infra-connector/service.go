package infra_connector

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

const (
	EndPoint = "infra-connector"
)

var (
	defaultPage = &model.PageQuery{
		Index: 0,
		Size:  100,
	}
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

func (s *Service) GetCloudAccounts(query *MdsCloudAccountsQuery) (model.Paged[model.MdsByocCloudAccount], error) {
	var response model.Paged[model.MdsByocCloudAccount]
	if query == nil {
		return response, fmt.Errorf("query cannot be nil")
	}

	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, CloudAccount)

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Service) GetCertificates(query *MDSCertificateQuery) (model.Paged[model.MdsByocCertificate], error) {
	var response model.Paged[model.MdsByocCertificate]
	if query == nil {
		return response, fmt.Errorf("query cannot be nil")
	}

	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, Certificate)

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Service) GetTshirtSizes(query *MdsTshirtSizesQuery) (model.Paged[model.MdsByocTshirtSize], error) {
	var response model.Paged[model.MdsByocTshirtSize]
	if query == nil {
		return response, fmt.Errorf("query cannot be nil")
	}

	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, K8sCluster, TshirtSize)

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Service) GetProviderTypes() ([]string, error) {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, CloudAccount, Types)
	var response []string

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (s *Service) GetCloudProviderRegions() ([]model.MdsByocCloudProviderRegion, error) {
	var response []model.MdsByocCloudProviderRegion

	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, CloudProviders)

	_, err := s.Api.Get(&reqUrl, nil, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
