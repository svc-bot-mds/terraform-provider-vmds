terraform {
  required_providers {
    mds = {
      source = "hashicorp.com/edu/mds"
    }
  }
}

provider "mds" {
  host      = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}


data "mds_policies" "policies" {
}

output "policies_data" {
  value = data.mds_policies.policies
}

resource "mds_service_account" "svc_account" {
  name = "test-svc-tf-update1"
  tags = ["update-svc-acct","from-tf"]
  policy_ids = ["64539a8f7d85190f7e5ae1e1","644a15e5df79724efa6b77ec"]
}