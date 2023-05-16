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

resource "vmds_policy" "policy_rabbitmq" {
  name = "test-from-tf-dont-use-2"
  service_type = "RABBITMQ"
  permission_specs = [
    {permissions: ["monitoring"], role: "monitoring", resource: "cluster:audit-test-11"}
  ]
}