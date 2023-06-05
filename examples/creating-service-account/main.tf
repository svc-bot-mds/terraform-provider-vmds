terraform {
  required_providers {
    vmds = {
      source = "svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host      = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

locals {
  policies = ["test-svc-pol"]
}

data "vmds_policies" "all" {
}

output "policies_data" {
  value = data.vmds_policies.all
}

resource "vmds_service_account" "test" {
  name       = "test-svc-tf-testing-131"
  tags       = ["update-svc-acct", "from-tf"]
  policy_ids = [for policy in data.vmds_policies.all.policies : policy.id if contains(local.policies, policy.name)]

  //Oauth app details
  oauth_app = {
    description = " description1"
    ttl_spec    = {
      ttl       = "1"
      time_unit = "HOURS"
    }
  }
}