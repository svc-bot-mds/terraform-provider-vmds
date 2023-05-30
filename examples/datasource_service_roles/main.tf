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

data "vmds_service_roles" "roles"{
  type = "RABBITMQ"
}

output "service_roles" {
  value = data.vmds_service_roles.roles
}
