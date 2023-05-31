// pass valid data with respect to the instance type selected
data "vmds_regions" "dedicated_aws" {
  cpu                  = "1"
  cloud_provider       = "aws"
  memory               = "4Gi"
  storage              = "4Gi"
  node_count           = "1"
  dedicated_data_plane = true
}