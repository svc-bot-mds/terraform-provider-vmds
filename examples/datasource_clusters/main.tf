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

data "vmds_clusters" "cluster_list"{
  service_type = "RABBITMQ"
}

output "clusters_data" {
  value = data.vmds_clusters.cluster_list
}
