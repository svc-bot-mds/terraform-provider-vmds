package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
)

const (
	EndPoint = "authservice"
)

type Service struct {
	*core.Service
}

func NewService(hostUrl *string, root *core.Root) *Service {
	return &Service{
		Service: core.NewService(hostUrl, EndPoint, root),
	}
}

// GetAccessToken - Get a new token for user
func (s *Service) GetAccessToken() (*TokenResponse, error) {
	if s.Api.AuthToUse.ApiToken == "" {
		return nil, fmt.Errorf("define API Token")
	}

	reqUrl := fmt.Sprintf("%s/%s", s.Endpoint, Token)

	tokenRequest := TokenRequest{
		ApiKey:        s.Api.AuthToUse.ApiToken,
		ClientId:      s.Api.AuthToUse.ClientId,
		ClientSecret:  s.Api.AuthToUse.ClientSecret,
		AccessToken:   s.Api.AuthToUse.AccessToken,
		OAuthAppTypes: s.Api.AuthToUse.OAuthAppType,
		OrgId:         s.Api.AuthToUse.OrgId,
	}
	body, err := s.Api.Post(&reqUrl, &tokenRequest, nil)
	if err != nil {
		return nil, err
	}

	ar := TokenResponse{
		Token: string(body),
	}

	err = s.processAuthResponse(&ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (s *Service) processAuthResponse(response *TokenResponse) error {
	s.Api.Token = &response.Token
	token, err := jwt.Parse(*s.Api.Token, nil)
	if token == nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	s.Api.OrgId = claims["context_name"].(string)
	return nil
}
