package main

import (
	"errors"
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/oauth_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/customer-metadata"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/model"
)

func main() {
	mdsHost := "MDS_HOST_URL"
	client, err := mds.NewClient(&mdsHost, &model.ClientAuth{
		ApiToken:     "API_TOKEN",
		OAuthAppType: oauth_type.ApiToken,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.CustomerMetadata.UpdateMdsUser("64533d8a2cee5b76e7c5fa70", &customer_metadata.MdsUserUpdateRequest{
		//PolicyIds:   []string{"644a14ac4efa951adae6b7d3"},
		Tags: []string{"client-test"},
		ServiceRoles: &[]customer_metadata.RolesRequest{
			{RoleId: "ManagedDataService:Admin"},
		},
	})

	fmt.Println(err)
	apiErr := core.ApiError{}
	fmt.Println(err != nil && errors.As(err, &apiErr), apiErr)
}
