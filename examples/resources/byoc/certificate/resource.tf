resource "vmds_certificate" "example" {
  name          = "tf-certificate-12"
  provider_type = "<<PROVIDER TYPE>>"
  domain_name    = "<<DOMAIN_NAME>>"
  certificate = "<<CERTIFICATE>>"
  certificate_ca = "<<CERTIFICATE_CA>>"
  certificate_key = "<<CERTIFICATE_KEY>>"

  //non editable fields during the update
  lifecycle {
    ignore_changes = [name, provider_type, domain_name]
  }
}
