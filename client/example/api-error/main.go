package main

import (
	"errors"
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/oauth_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds/core"
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

	_, err = client.Controller.GetMdsCluster("12376yhsjdasd")

	if err != nil {
		fmt.Println(err)
		var apiError core.ApiError
		if errors.As(err, &apiError) {
			fmt.Println("recognized")
			fmt.Println(apiError.ErrorMessage)
		}
		return
	}
}
