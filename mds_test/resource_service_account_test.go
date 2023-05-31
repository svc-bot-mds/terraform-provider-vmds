package mds_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* Set up any prerequisites or check for required dependencies */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `locals {
  											account_type  = "SERVICE_ACCOUNT"
  											policies = ["gya-policy","eu301"]
										}

										data "vmds_policies" "policies" {
										}

										output "policies_data" {
  											value = data.vmds_policies.policies
										}

										resource "vmds_service_account" "svc_account" {
  											name = "test-svc-tf-create-sa"
  											tags = ["create-svc-acct","from-tf"]
  											policy_ids =  [for policy in data.vmds_policies.policies.policies: policy.id if contains(local.policies, policy.name) ]

  											// non editable fields
  											lifecycle {
   											 ignore_changes = [name]
  											}
									}

data "vmds_service_accounts" "service_accounts" {
}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_service_accounts.service_accounts", "id"),
					resource.TestCheckResourceAttr("data.vmds_service_accounts.service_accounts", "service_accounts.0.name", "test-svc-tf-create-sa"),
				),
			},
			{
				Config: providerConfig + `locals {
  											account_type  = "SERVICE_ACCOUNT"
										}

										resource "vmds_service_account" "svc_account" {
  											name = "test-svc-tf-create-sa"
  											tags = ["update-svc-acct"]
  											policy_ids =  [for policy in data.vmds_policies.policies.policies: policy.id if contains(local.policies, policy.name) ]

  											// non editable fields
  											lifecycle {
   											 ignore_changes = [name]
  											}
									}
data "vmds_service_accounts" "service_accounts" {
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmds_service_accounts.service_accounts", "id"),
					resource.TestCheckResourceAttr("data.vmds_service_accounts.service_accounts", "service_accounts.0.name", "test-svc-tf-create-sa"),
				),
			},
			{
				Config: providerConfig,
			},
		},
	})
}
