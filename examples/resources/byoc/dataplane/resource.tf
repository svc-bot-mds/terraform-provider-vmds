
data vmds_certificates "all" {

}

output "certificate_list" {
  value = data.vmds_certificates.all
}
data vmds_cloud_accounts "all" {

}

output "cloud_accounts" {
  value = data.vmds_cloud_accounts.all
}
data vmds_cloud_provider_regions "all" {

}

output "res" {
  value = data.vmds_cloud_provider_regions.all
}
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