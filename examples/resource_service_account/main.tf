terraform {
  required_providers {
    vmds = {
      source = "hashicorp.com/edu/vmds"
    }
  }
}

provider "vmds" {
  host     = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

locals {
  policies = ["gya-policy","eu301"]
}

data "vmds_policies" "policies" {
}

output "policies_data" {
  value = data.vmds_policies.policies
}

resource "vmds_service_account" "svc_account" {
  name = "test-svc-tf-update1"
  tags = ["update-svc-acct","from-tf"]
  policy_ids =  [for policy in data.vmds_policies.policies.policies: policy.id if contains(local.policies, policy.name) ]

  // non editable fields
  lifecycle {
    ignore_changes = [name]
  }
}