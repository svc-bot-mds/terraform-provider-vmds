terraform {
  required_providers {
    vmds = {
      source = "hashicorp.com/edu/vmds"
    }
  }
}

provider "vmds" {
  host     = "MDS_HOST"
  api_token = "API_TOKEN"
}

locals {
  account_type  = "USER_ACCOUNT"
  service_roles = ["Developer", "Admin"]
}

data "vmds_roles" "all" {
}

output "roles_data" {
  value = data.vmds_roles.all
}

data "vmds_policies" "policies" {
}

output "policies_data" {
  value = data.vmds_policies.policies
}

resource "vmds_user" "temp" {
  email      = "developer@vmware.com"
  policy_ids = ["64539a8f7d85190f7e5ae1e1"]
  tags       = ["new-user-tf", "update-tf-user"]
  role_ids   = [for role in data.vmds_roles.all.roles : role.role_id if contains(local.service_roles, role.name)]
  timeouts   = {
    create = "1m"
  }
}
