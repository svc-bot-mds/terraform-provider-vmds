// example of a RABBITMQ policy
resource "vmds_policy" "rabbitmq" {
  name             = "test-tf"
  service_type     = "RABBITMQ"
  permission_specs = [
    {
      permissions = ["read"],
      role        = "read",
      resource    = "cluster:example-cluster-name"
    },
    {
      "permissions" = ["write"],
      "role"        = "write",
      "resource"    = "cluster:example-cluster-name/queue:my-queue"
    }
  ]
}

// example of a NETWORK policy
resource "vmds_policy" "network" {
  name         = "network-policy-from-tf"
  service_type = "NETWORK"
  network_spec = {
    cidr             = "10.22.55.0/24",
    network_port_ids = ["rmq-streams", "rmq-amqps"]
  }
}