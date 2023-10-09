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
  type = string
  default = <<EOF
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
  type = string
  default = <<EOF
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

resource "vmds_cloud_account" "example" {
  name          = "tf-cloud-account1"
  provider_type = "vrli"
  credential    = var.vrli_cred

  //non editable fields during the update
  lifecycle {
    ignore_changes = [name]
  }
}