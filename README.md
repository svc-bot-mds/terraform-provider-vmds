# Terraform Provider for VMware Managed Data Services

## About

This repository contains code for the Terraform Provider for VMware Managed Data Services. It supports provisioning of Clusters/Instances of Services (currently only RabbitMQ) and access management of Users on those resources.

## Configuration

The Terraform Provider for VMware MDS is available via the Terraform Registry: [svc-ops-mds/vmds](https://registry.terraform.io/providers/svc-bot-mds/vmds). To be able to use it successfully, please use below snippet to set up the provider:

```hcl
terraform {
  required_providers {
    vmds = {
      source = "svc-bot-mds/vmds"
    }
  }
}

provider "vmds" {
  host      = "https://mds-console.example.com" # (required) the URL of hosted MDS
  api_token = "XXXXXX__API_TOKEN__XXXXXX"       # (required) can be generated from CSP > Accounts page
}
```
