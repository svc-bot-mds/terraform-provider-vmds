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

locals {
  service_type       = "RABBITMQ"
  provider           = "aws"
  policy_with_create = ["open-to-all"]
  policy_with_update = ["custom-nw"]
  instance_type      = "XX-SMALL"
}

data "vmds_regions" "all" {
  cloud_provider = "aws"
  instance_size = "XX-SMALL"
}

output "regions" {
  value = data.vmds_regions.all
}

data "vmds_network_policies" "create" {
  names = local.policy_with_create
}

data "vmds_network_policies" "update" {
  names = local.policy_with_update
}

output "network_policies_data" {
  value = {
    update = data.vmds_network_policies.update
    create = data.vmds_network_policies.create
  }
}

resource "vmds_cluster" "test" {
  name               = "my-rmq-cls"
  service_type       = local.service_type
  cloud_provider     = local.provider
  instance_size      = local.instance_type
  region             = data.vmds_regions.all.regions[0].id
  network_policy_ids = data.vmds_network_policies.create.policies[*].id
  tags               = ["mds-tf", "example", "new-tag"]
  timeouts           = {
    create = "10m"
  }
}


resource "vmds_cluster_network_policies_association" "test" {
  id         = vmds_cluster.test.id
  policy_ids = data.vmds_network_policies.update.policies[*].id
}