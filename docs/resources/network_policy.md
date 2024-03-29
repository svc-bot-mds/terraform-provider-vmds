---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vmds_network_policy Resource - vmds"
subcategory: ""
description: |-
  Represents a policy on MDS.
---

# vmds_network_policy (Resource)

Represents a policy on MDS.

## Example Usage

```terraform
// example of a NETWORK policy
resource "vmds_network_policy" "network" {
  name         = "network-policy-from-tf"
  network_spec = {
    cidr             = "10.22.55.0/24",
    network_port_ids = ["rmq-streams", "rmq-amqps"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the policy
- `network_spec` (Attributes) Network config to allow access to service resource. (see [below for nested schema](#nestedatt--network_spec))
- `service_type` (String) Type of policy to manage. Supported values is:  `NETWORK`.

### Read-Only

- `id` (String) Auto-generated ID of the policy after creation, and can be used to import it from MDS to terraform state.
- `resource_ids` (Set of String) IDs of service resources/instances being managed by the policy.

<a id="nestedatt--network_spec"></a>
### Nested Schema for `network_spec`

Required:

- `cidr` (String) CIDR value to allow access from. Ex: `10.45.66.80/30`
- `network_port_ids` (Set of String) IDs of network ports to open up for access. Please make use of datasource `vmds_network_ports` to get IDs of ports available for services.

## Import

Import is supported using the following syntax:

```shell
# Policy can be imported by specifying the alphanumeric identifier.
terraform import vmds_policy.network s546dg29fh2ksh3dfr
```
