terraform {
  required_providers {
    vmds = {
      source = "vmware/managed-data-services"
    }
  }
}

provider "vmds" {
  host     = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

data "vmds_cluster_metadata_by_id" "metadata1" {
id = "644a152bbaa9ff65a87bd139"
}

output "network_ports_data" {
  value = data.vmds_cluster_metadata_by_id.metadata1
}

