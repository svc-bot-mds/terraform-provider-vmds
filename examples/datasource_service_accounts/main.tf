terraform {
  required_providers {
    vmds = {
      source = "vmware/managed-data-services"
    }
  }
}

provider "vmds" {
  host      = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

data "vmds_service_accounts" "service_accounts" {
}

output "service_accounts_data" {
  value = data.vmds_service_accounts.service_accounts
}