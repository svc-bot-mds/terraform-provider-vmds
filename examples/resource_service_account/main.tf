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


data "vmds_policies" "policies" {
}

output "policies_data" {
  value = data.vmds_policies.policies
}

resource "vmds_service_account" "svc_account" {
  name = "test-svc-tf-update1"
  tags = ["update-svc-acct","from-tf"]
  policy_ids = ["64539a8f7d85190f7e5ae1e1","644a15e5df79724efa6b77ec"]
}