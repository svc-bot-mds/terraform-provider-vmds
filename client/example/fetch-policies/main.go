package main

import (
	"fmt"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/oauth_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/constants/policy_type"
	"github.com/svc-bot-mds/terraform-provider-vmds/client/mds"
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

	response, err := client.CustomerMetadata.GetPolicies(&customer_metadata.MdsPoliciesQuery{
		Type:  policy_type.NETWORK,
		Names: []string{"my-nw-policy"},
	})

	fmt.Println(response.Get())
	for _, dto := range *response.Get() {
		fmt.Println(dto)
	}
}
