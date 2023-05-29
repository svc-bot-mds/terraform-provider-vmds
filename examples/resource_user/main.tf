terraform {
  required_providers {
    vmds = {
      source = "vmware/managed-data-services"
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
  policies = ["gya-policy","eu301"]
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
  email      = "developer11@vmware.com"
  tags       = ["new-user-tf", "update-tf-user"]
  role_ids   = [for role in data.vmds_roles.all.roles : role.role_id if contains(local.service_roles, role.name)]
  policy_ids = [for policy in data.vmds_policies.policies.policies: policy.id if contains(local.policies, policy.name) ]
  timeouts   = {
    create = "1m"
  }

  // non editable fields
  lifecycle {
    ignore_changes = [email, status]
  }
}
