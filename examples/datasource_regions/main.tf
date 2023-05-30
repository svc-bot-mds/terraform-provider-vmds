terraform {
  required_providers {
    vmds = {
      source = "svc-bot-mds/vmds"
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

// pass valid data with respect to the instance type selected
data "vmds_regions" "available_regions" {
  cpu                  = "1"
  cloud_provider       = local.provider
  memory               = "4Gi"
  storage              = "4Gi"
  node_count           = "1"
  dedicated_data_plane = true
}

output "regions_data" {
  value = data.vmds_regions.available_regions
}
