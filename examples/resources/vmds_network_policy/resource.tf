
// example of a NETWORK policy
resource "vmds_network_policy" "network" {
  name         = "network-policy-from-tf"
  network_spec = {
    cidr             = "10.22.55.0/24",
    network_port_ids = ["rmq-streams", "rmq-amqps"]
  }
}