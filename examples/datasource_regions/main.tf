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
  provider = "aws"
}

data "vmds_regions" "all" {
  cpu                  = "1"
  cloud_provider       = local.provider
  memory               = "4Gi"
  storage              = "4Gi"
  node_count           = "1"
  dedicated_data_plane = false
}

output "regions_data" {
  value = data.mds_regions.all
}
