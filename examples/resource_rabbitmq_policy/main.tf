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
  name = "test-svc-policy-from-tf"
  service_type = "RABBITMQ"
  permission_spec = [
    {permissions: ["monitoring"], role: "monitoring", resource: "cluster:audit-test-11"}
  ]
  network_specs = []
}