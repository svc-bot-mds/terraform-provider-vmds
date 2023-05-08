terraform {
  required_providers {
    vmds = {
      source = "hashicorp.com/edu/vmds"
    }
  }
}

provider "vmds" {
  host      = "MDS_HOST_URL"
  api_token = "API_TOKEN"
}

locals {
  service_type       = "RABBITMQ"
  provider           = "aws"
  policy_with_create = ["open-to-all"]
  policy_with_update = ["custom-nw"]
}

data "vmds_regions" "all" {
  cpu                  = "1"
  cloud_provider       = local.provider
  memory               = "4Gi"
  storage              = "4Gi"
  node_count           = "1"
  dedicated_data_plane = false
}

data "vmds_network_policies" "create" {
  names = local.policy_with_create
}

data "vmds_network_policies" "update" {
  names = local.policy_with_update
}

output "network_policies_data" {
  value = {
    update = data.mds_network_policies.update
    create = data.mds_network_policies.create
  }
}

resource "vmds_cluster" "test" {
  name               = "my-rmq-cls"
  service_type       = local.service_type
  cloud_provider     = local.provider
  instance_size      = local.instance_type
  region             = local.region
  network_policy_ids = data.mds_network_policies.create.policies[*].id
  tags               = ["mds-tf", "example", "new-tag"]
  timeouts           = {
    create = "10m"
  }
}


resource "vmds_cluster_network_policies_association" "test" {
  id         = mds_cluster.test.id
  policy_ids = data.mds_network_policies.update.policies[*].id
}