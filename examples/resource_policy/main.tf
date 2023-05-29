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

locals {
  cluster_name ="test"
  queue_name ="dc"
}

data "vmds_clusters" "cluster_list"{
  service_type = "RABBITMQ"
}

data "vmds_service_roles" "roles"{
  type = "RABBITMQ"
}

data "vmds_cluster_metadata" "metadata" {
  id = "6465f3ae265b393b4e42e9bd"
}

output "cluster_metadata" {
  value = data.vmds_cluster_metadata.metadata
}

resource "vmds_policy" "policy_rabbitmq" {
  name = "test-tf"
  service_type = "RABBITMQ"
  permission_specs = [
    {
      permissions: ["read"],
      role: "read",
      resource: "cluster:${local.cluster_name}"
    },
    {
      "permissions":["write"],
      "role":"write",
      "resource":"cluster:${local.cluster_name}/queue:${local.queue_name}"}
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