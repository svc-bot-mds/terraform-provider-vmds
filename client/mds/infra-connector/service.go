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

func (s *Service) GetCloudAccounts(query *MdsCloudAccountsQuery) (model.Paged[model.MdsCloudAccount], error) {
	var response model.Paged[model.MdsCloudAccount]
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

func (s *Service) GetCertificates(query *MDSCertificateQuery) (model.Paged[model.MdsCertificate], error) {
	var response model.Paged[model.MdsCertificate]
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

func (s *Service) GetTshirtSizes(query *MdsTshirtSizesQuery) (model.Paged[model.MdsTshirtSize], error) {
	var response model.Paged[model.MdsTshirtSize]
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

func (s *Service) GetDataPlaneRegions() ([]model.MdsDataPlaneRegion, error) {
	var response []model.MdsDataPlaneRegion

	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, CloudProviders)

	_, err := s.Api.Get(&reqUrl, nil, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// CreateDataPlane - Submits a request to create dataplane
func (s *Service) CreateDataPlane(requestBody *DataPlaneCreateRequest) (*model.TaskResponse, error) {
	if requestBody == nil {
		return nil, fmt.Errorf("requestBody cannot be nil")
	}
	var response model.TaskResponse
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, K8sCluster)

	_, err := s.Api.Post(&urlPath, requestBody, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

func (s *Service) GetDataPlanes(query *DataPlaneQuery) (model.Paged[model.DataPlane], error) {
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, K8sCluster)
	var response model.Paged[model.DataPlane]

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&urlPath, query, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *Service) GetDataPlaneById(id string) (model.DataPlane, error) {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, K8sCluster, id)
	var response model.DataPlane

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// DeleteDataPlane - Submits a request to delete dataplane
func (s *Service) DeleteDataPlane(id string) error {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, K8sCluster, id)

	_, err := s.Api.Delete(&urlPath, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
