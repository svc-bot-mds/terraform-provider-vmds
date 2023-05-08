package controller

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/utils"
	"strings"
)

const (
	EndPoint = "controller"
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

// GetMdsClusters - Returns page of clusters
func (s *Service) GetMdsClusters(query *MdsClustersQuery) (model.Paged[model.MdsCluster], error) {
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Clusters)
	var response model.Paged[model.MdsCluster]

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&urlPath, query, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// GetAllMdsClusters - Returns list of all clusters
func (s *Service) GetAllMdsClusters(query *MdsClustersQuery) ([]model.MdsCluster, error) {
	var clusters []model.MdsCluster
	for {
		queriedClusters, err := s.GetMdsClusters(query)
		if err != nil {
			return clusters, err
		}
		clusters = append(clusters, *queriedClusters.Get()...)
		nextPage := utils.GetNextPageInfo(queriedClusters.GetPage())
		if nextPage == nil {
			break
		}
		query.PageQuery = *nextPage
	}
	return clusters, nil
}

// GetMdsCluster - Returns the cluster by ID
func (s *Service) GetMdsCluster(id string) (*model.MdsCluster, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Clusters, id)
	var response model.MdsCluster

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// CreateMdsCluster - Submits a request to create cluster
func (s *Service) CreateMdsCluster(requestBody *MdsClusterCreateRequest) (*model.TaskResponse, error) {
	if requestBody == nil {
		return nil, fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Clusters)
	var response model.TaskResponse

	_, err := s.Api.Post(&urlPath, requestBody, &response)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// UpdateMdsCluster - Submits a request to update cluster
func (s *Service) UpdateMdsCluster(id string, requestBody *MdsClusterUpdateRequest) (*model.MdsCluster, error) {
	if id == "" {
		return nil, fmt.Errorf("cluster ID cannot be empty")
	}
	if requestBody == nil {
		return nil, fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Clusters, id)
	var response model.MdsCluster

	_, err := s.Api.Patch(&urlPath, requestBody.Tags, &response)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// UpdateMdsClusterNetworkPolicies - Submits a request to update cluster network policies
func (s *Service) UpdateMdsClusterNetworkPolicies(id string, requestBody *MdsClusterNetworkPoliciesUpdateRequest) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("cluster ID cannot be empty")
	}
	if requestBody == nil {
		return nil, fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s/%s/%s", s.Endpoint, Clusters, id, NetworkPolicy)

	bodyBytes, err := s.Api.Patch(&urlPath, requestBody, nil)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

// DeleteMdsCluster - Submits a request to delete cluster
func (s *Service) DeleteMdsCluster(id string) (*model.TaskResponse, error) {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Clusters, id)
	var response model.TaskResponse

	_, err := s.Api.Delete(&urlPath, nil, &response)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// GetServiceInstanceTypes - Returns list of clusters
func (s *Service) GetServiceInstanceTypes(serviceTypeQuery *MdsInstanceTypesQuery) (model.MdsInstanceTypeList, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s", s.Endpoint, Services, InstanceTypes)
	var response model.MdsInstanceTypeList

	_, err := s.Api.Get(&reqUrl, serviceTypeQuery, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
