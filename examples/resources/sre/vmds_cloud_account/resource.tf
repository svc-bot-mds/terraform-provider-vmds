terraform {
  required_providers {
    vmds = {
      source = "svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host = "MDS_HOST_URL"

  //Get the authentication with "username and password"
  username = "USERNAME"
  password = "PASSWORD"
  org_id   = "ORG-ID"
  type     = "user_creds"
}

variable "vrli_cred" {
  description = "VRLI CRED JSON"
  type        = string
  default     = <<EOF
  {
    "apiEndpoint":"test11",
    "apiKey":"fgeddsd"
  }
EOF
}

variable "aws_cred" {
  description = "AWS CRED JSON"
  type        = string
  default     = <<EOF
{
    "ACCESS_KEY_ID": "REPLACE_ACCESS_KEY_ID",
    "SECRET_ACCESS_KEY" : "REPLACE_SECRET_ACCESS_KEY",
    "targetAmazonAccountId" : "REPLACE_AWS_ACCOUNT",
    "powerUserRoleNameInTargetAccount" : "REPLACE_AWS_SERVICE_ACCOUNT_ROLE"
  }
EOF
}

variable "gcp_cred" {
  description = "GCP CRED JSON"
  type        = string
  default     = <<EOF
{
    "type": "test",
    "project_id" : "test",
    "private_key_id" : "test",
    "private_key" : "test",
    "client_email" : "test",
    "client_id": "test",
    "auth_uri" : "test",
    "token_uri" : "test",
    "auth_provider_x509_cert_url" : "test",
    "client_x509_cert_url" : "test"
  }
EOF
}

variable "tkgs_cred" {
  description = "TKGs CRED JSON"
  type        = string
  default     = <<EOF
{
    "userName": "test",
    "password": "REPLACE",
    "supervisorManagementIP": "SOME_IP",
    "vsphereNamespace": "NAMESPACE"
}
EOF
}

data "vmds_provider_types" "create" {
}

output "provider_types" {
  value = {
    create = data.vmds_provider_types.create
  }
}
resource "vmds_cloud_account" "example" {
  name          = "tf-cloud-account1"
  provider_type = element(data.vmds_provider_types.create.list, 0)
  credentials   = var.tkgs_cred
  shared        = true
  tags          = ["tag1", "tag2"]

  //non editable fields during the update
  lifecycle {
    ignore_changes = [name]
  }
}

