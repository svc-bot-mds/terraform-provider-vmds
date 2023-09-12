data "vmds_certificates" "all" {
}

output "resp" {
  value = data.vmds_certificates.all
}