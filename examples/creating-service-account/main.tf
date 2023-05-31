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
  policies = ["viewer-policy", "eu301"]
}

data "vmds_policies" "all" {
}

output "policies_data" {
  value = data.vmds_policies.all
}

resource "vmds_service_account" "test" {
  name       = "test-svc-tf-update1"
  tags       = ["update-svc-acct", "from-tf"]
  policy_ids = [for policy in data.vmds_policies.all.policies : policy.id if contains(local.policies, policy.name)]

  // non editable fields
  lifecycle {
    ignore_changes = [name]
  }
}