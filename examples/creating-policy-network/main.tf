terraform {
  required_providers {
    vmds = {
      source = "svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host      = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

// network port IDs can be referred using this datasource
data "vmds_network_ports" "all" {
}

output "cluster_metadata" {
  value = data.vmds_network_ports.all
}

resource "vmds_policy" "network" {
  name         = "network-policy-from-tf"
  service_type = "NETWORK"
  network_spec = {
    cidr             = "10.22.55.0/24",
    network_port_ids = ["rmq-streams", "rmq-amqps"]
  }
}