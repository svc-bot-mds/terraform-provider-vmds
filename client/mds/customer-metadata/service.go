package customer_metadata

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/account_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"strings"
)

const (
	EndPoint = "customermetadata"
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

// GetPolicies - Returns list of Policies
func (s *Service) GetPolicies(query *MdsPoliciesQuery) (model.Paged[model.MDSPolicies], error) {
	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, Policies)
	var response model.Paged[model.MDSPolicies]

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetMdsUsers - Return list of Users
func (s *Service) GetMdsUsers(query *MdsUsersQuery) (model.Paged[model.MdsUser], error) {
	var response model.Paged[model.MdsUser]
	if query == nil {
		return response, fmt.Errorf("query cannot be nil")
	}

	query.AccountType = account_type.USER_ACCOUNT
	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, Users)

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// CreateMdsUser - Submits a request to create user
func (s *Service) CreateMdsUser(requestBody *MdsCreateUserRequest) error {
	if requestBody == nil {
		return fmt.Errorf("requestBody cannot be nil")
	}
	requestBody.AccountType = account_type.USER_ACCOUNT
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Users)

	_, err := s.Api.Post(&urlPath, requestBody, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMdsUser - Submits a request to update user
func (s *Service) UpdateMdsUser(id string, requestBody *MdsUserUpdateRequest) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if requestBody == nil {
		return fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)

	_, err := s.Api.Patch(&urlPath, requestBody, nil)
	return err
}

// GetMdsUser - Returns the user by ID
func (s *Service) GetMdsUser(id string) (*model.MdsUser, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)
	var response model.MdsUser

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// DeleteMdsUser - Submits a request to delete user
func (s *Service) DeleteMdsUser(id string) error {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)

	_, err := s.Api.Delete(&urlPath, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetMdsServiceAccounts - Return list of Service Accounts
func (s *Service) GetMdsServiceAccounts(query *MdsServiceAccountsQuery) (model.Paged[model.MdsServiceAccount], error) {

	var response model.Paged[model.MdsServiceAccount]
	if query == nil {
		return response, fmt.Errorf("query cannot be nil")
	}

	query.AccountType = account_type.SERVICE_ACCOUNT
	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, Users)

	if query.Size == 0 {
		query.Size = defaultPage.Size
	}

	_, err := s.Api.Get(&reqUrl, query, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// CreateMdsServiceAccount - Submits a request to create service account
func (s *Service) CreateMdsServiceAccount(requestBody *MdsCreateSvcAccountRequest) error {
	if requestBody == nil {
		return fmt.Errorf("requestBody cannot be nil")
	}

	requestBody.AccountType = account_type.SERVICE_ACCOUNT

	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Users)

	_, err := s.Api.Post(&urlPath, requestBody, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMdsServiceAccount - Submits a request to update service account
func (s *Service) UpdateMdsServiceAccount(id string, requestBody *MdsSvcAccountUpdateRequest) error {
	if id == "" {
		return fmt.Errorf("service account ID cannot be empty")
	}
	if requestBody == nil {
		return fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)

	_, err := s.Api.Patch(&urlPath, requestBody, nil)
	return err
}

// GetMdsServiceAccount - Returns the service account by ID
func (s *Service) GetMdsServiceAccount(id string) (*model.MdsServiceAccount, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)
	var response model.MdsServiceAccount

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// DeleteMdsServiceAccount - Submits a request to delete service account
func (s *Service) DeleteMdsServiceAccount(id string) error {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Users, id)

	_, err := s.Api.Delete(&urlPath, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreatePolicy - Submits a request to create policy
func (s *Service) CreatePolicy(requestBody *MdsCreateUpdatePolicyRequest) (*model.MDSPolicies, error) {
	if requestBody == nil {
		return nil, fmt.Errorf("requestBody cannot be nil")
	}
	var response model.MDSPolicies
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Policies)

	_, err := s.Api.Post(&urlPath, requestBody, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// UpdateMdsPolicy - Submits a request to update policy
func (s *Service) UpdateMdsPolicy(id string, requestBody *MdsCreateUpdatePolicyRequest) error {
	if id == "" {
		return fmt.Errorf("policy ID cannot be empty")
	}
	if requestBody == nil {
		return fmt.Errorf("requestBody cannot be nil")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Policies, id)

	_, err := s.Api.Put(&urlPath, requestBody, nil)
	return err
}

// GetMDSPolicy - Submits a request to fetch policy
func (s *Service) GetMDSPolicy(id string) (*model.MDSPolicies, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Policies, id)
	var response model.MDSPolicies

	_, err := s.Api.Get(&urlPath, nil, &response)
	if err != nil {
		return &response, err
	}

	return &response, err
}

// DeleteMdsPolicy - Submits a request to delete policy
func (s *Service) DeleteMdsPolicy(id string) error {
	urlPath := fmt.Sprintf("%s/%s/%s", s.Endpoint, Policies, id)

	_, err := s.Api.Delete(&urlPath, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
