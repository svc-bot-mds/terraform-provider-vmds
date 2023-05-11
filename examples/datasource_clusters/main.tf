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

data "vmds_clusters" "cluster_list"{
  service_type = "RABBITMQ"
}

output "clusters_data" {
  value = data.vmds_clusters.cluster_list
}
