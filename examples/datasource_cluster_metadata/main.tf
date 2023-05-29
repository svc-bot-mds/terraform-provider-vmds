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

  data "vmds_cluster_metadata" "metadata" {
    id = "CLuster_Id"
  }

output "network_ports_data" {
  value = data.vmds_cluster_metadata.metadata
}

