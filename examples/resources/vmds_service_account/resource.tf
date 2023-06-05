resource "vmds_service_account" "example" {
  name       = "example-acc"
  tags       = ["temporary", "limited-role"]
  policy_ids = ["as73i83jnfkw9wr"]

  // non editable fields
  lifecycle {
    ignore_changes = [name]
  }
  //Oauth app details
  oauth_app = {
    description = "Oauth app created for example-acc service account"
    ttl_spec    = {
      ttl       = "1"
      time_unit = "HOURS"
    }
  }
}