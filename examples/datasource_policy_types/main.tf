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

data "vmds_policy_types" "typesList" {
}

output "network_ports_data" {
  value = data.vmds_policy_types.typesList
}

