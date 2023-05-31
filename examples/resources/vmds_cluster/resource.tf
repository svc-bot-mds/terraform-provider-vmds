resource "vmds_cluster" "example" {
  name               = "example-rmq-cls"
  service_type       = "RABBITMQ"
  cloud_provider     = "aws"
  instance_size      = "XX-SMALL"
  region             = "us-east-1"
  network_policy_ids = ["ajgynfg634bfj63hd"]
  tags               = ["mds-tf", "example"]
  timeouts           = {
    create = "3m"
    delete = "1m"
  }
  // non editable fields
  lifecycle {
    ignore_changes = [instance_size, name, cloud_provider, region, service_type]
  }
}