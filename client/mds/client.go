package mds

import (
	"crypto/tls"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/auth"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/controller"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/infra-connector"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/service-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/upgrade-service"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default MDS URL
const HostURL string = "http://localhost:8080"

// Client -
type Client struct {
	Root             *core.Root
	Auth             *auth.Service
	Controller       *controller.Service
	InfraConnector   *infra_connector.Service
	CustomerMetadata *customer_metadata.Service
	ServiceMetadata  *service_metadata.Service
	UpgradeService   *upgrade_service.Service
}

// NewClient -
func NewClient(host *string, authInfo *model.ClientAuth) (*Client, error) {
	hostUrl := HostURL
	if len(strings.TrimSpace(*host)) != 0 {
		hostUrl = *host
	}

	httpClient := prepareHttpClient()
	root := &core.Root{
		// Default MDS URL
		HostUrl:    &hostUrl,
		AuthToUse:  authInfo,
		HttpClient: httpClient,
	}

	c := prepareClient(host, root)

	_, err := c.Auth.GetAccessToken()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func prepareClient(host *string, root *core.Root) *Client {
	return &Client{
		Root:             root,
		Auth:             auth.NewService(host, root),
		Controller:       controller.NewService(host, root),
		InfraConnector:   infra_connector.NewService(host, root),
		CustomerMetadata: customer_metadata.NewService(host, root),
		ServiceMetadata:  service_metadata.NewService(host, root),
		UpgradeService:   upgrade_service.NewService(host, root),
	}
}

func prepareHttpClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
