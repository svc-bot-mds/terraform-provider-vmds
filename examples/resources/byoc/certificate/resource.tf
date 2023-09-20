resource "vmds_certificate" "example" {
  name          = "tf-certificate-12"
  provider_type = "aws"
  domain_name    = "<<domain name>>"
  certificate_ca = "<<certificate ca>>"
  certificate = "<<certificate>>"
  certificate_key = "<certificate privte key>>"
}
