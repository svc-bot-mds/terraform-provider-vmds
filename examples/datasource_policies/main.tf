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
  value = data.mds_policies.policies
}

