terraform {
  required_providers {
    vmds = {
      source = "hashicorp.com/svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host      = "https://console.mds.vmware.com"

  username = " < Username > "
  password = " < Password > "

  type = "user_creds"
}

data "vmds_certificates" "all"{
}
output "resp" {
  value = data.vmds_certificates.all
}

