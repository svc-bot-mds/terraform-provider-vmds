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


data "vmds_network_policies" "network_policies" {
}

output "network_policies_data" {
  value = data.vmds_network_policies.network_policies
}

