package upgrade_service

import (
	"fmt"
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
	EndPoint = "upgradeservice"
)

type Service struct {
	*core.Service
}

func NewService(hostUrl *string, root *core.Root) *Service {
	return &Service{
		Service: core.NewService(hostUrl, EndPoint, root),
	}
}

// UpdateMdsClusterVersion updates the version of the MDS cluster
func (s *Service) UpdateMdsClusterVersion(id string, requestBody *UpdateMdsClusterVersionRequest) (*model.UpdateMdsClusterVersionResponse, error) {
	urlPath := fmt.Sprintf("%s/%s", s.Endpoint, Upgrade)
	var response model.UpdateMdsClusterVersionResponse

	fmt.Println("Version update Request : ", requestBody)
	fmt.Println("Version update Request : ", &requestBody)
	_, err := s.Api.Post(&urlPath, requestBody, &response)
	if err != nil {
		return &response, err
	}

	return &response, nil
}
