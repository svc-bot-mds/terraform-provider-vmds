resource "vmds_user" "example" {
  email      = "developer11@vmware.com"
  tags       = ["new-user", "viewer"]
  role_ids   = ["mds:viewer"]
  policy_ids = ["asdhh4bsd83bfd"]
  timeouts   = {
    create = "1m"
  }

  // non editable fields
  lifecycle {
    ignore_changes = [email]
  }
}