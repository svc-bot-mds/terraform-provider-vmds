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

resource "vmds_policy" "policy_network" {
  name = "network-policy-from-tf"
  service_type = "NETWORK"
  network_specs = [
    {cidr: "10.22.55.0/24", network_port_ids: ["rmq-metrics", "rmq-amqps"]}
  ]
  permission_spec = []
}