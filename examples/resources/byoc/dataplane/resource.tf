resource "vmds_byoc_dataplane" "example" {
  name    = "byoc-tf-test-1"
  account_id = "<<cloud account id>>"
  certificate_id = "<<certificate id>>"
  nodepool_type = "regular"
  region = "us-east-1"
  // non editable fields , edit is not allowed
  lifecycle {
    ignore_changes = [name, account_id, certificate_id, nodepool_type, region]
  }
}