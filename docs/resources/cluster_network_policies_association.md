---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vmds_cluster_network_policies_association Resource - vmds"
subcategory: ""
description: |-
  Represents the association between a service instance/cluster and NETWORK type policies.
---

# vmds_cluster_network_policies_association (Resource)

Represents the association between a service instance/cluster and `NETWORK` type policies.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) ID of the cluster.
- `policy_ids` (Set of String) IDs of the network policies to associate with the cluster.


