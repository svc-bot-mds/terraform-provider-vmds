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
  cluster_name ="test-1"
}

data "vmds_clusters" "cluster_list"{
  service_type = "RABBITMQ"
}

data "vmds_service_roles" "roles"{
  type = "RABBITMQ"
}

resource "vmds_policy" "policy_rabbitmq" {
  name = "test-tf-dont-use-2"
  service_type = "RABBITMQ"
  permission_specs = [
    {
      permissions: ["read"],
      role: "read",
      resource: "cluster:${local.cluster_name}"
    },
  ]
}

resource "vmds_policy" "policy_network" {
  name             = "network-policy-from-tf"
  service_type     = "NETWORK"
  network_spec     = {
    cidr : "10.22.55.0/24",
    network_port_ids : ["rmq-streams", "rmq-amqps"]
  }
}