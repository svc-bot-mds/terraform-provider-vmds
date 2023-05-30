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

locals {
  service_type = "RABBITMQ"
}

data "vmds_instance_types" "rmq" {
  service_type = local.service_type
}

output "instance_types" {
  value = data.vmds_instance_types.rmq
}