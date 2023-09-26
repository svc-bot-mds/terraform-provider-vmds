resource "vmds_cluster" "example" {
  name               = "test-terraform"
  cloud_provider = "aws"
  service_type       = "RABBITMQ"
  instance_size      = "XX-SMALL"
  region             = "eu-west-1"
  network_policy_ids = ["policy id"]
  tags               = ["mds-tf", "example"]
  dedicated = false
  shared = false

  // if cluster getting self hosted via byoc
  data_plane_id = "dataplane id"
  // non editable fields
  lifecycle {
    ignore_changes = [instance_size, name, cloud_provider, region, service_type]
  }
}