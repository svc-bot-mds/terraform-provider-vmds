terraform {
  required_providers {
    vmds = {
      source = "hashicorp.com/svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host      = "https://tdh-cp-vh.tdh.kr.com"

//  username = "venkatram.amalanathan@broadcom.com"
//  password = "Signin@07"
  username = "sre@broadcom.com"
  password = "VMware$123"
  type = "user_creds"
}

data "vmds_certificates" "all"{
}
output "resp" {
  value = data.vmds_certificates.all
}

