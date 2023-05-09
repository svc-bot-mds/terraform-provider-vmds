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

data "mds_service_accounts" "service_accounts" {
}

output "service_accounts_data" {
  value = data.mds_service_accounts.service_accounts
}