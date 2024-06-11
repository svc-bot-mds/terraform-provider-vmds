data "vmds_object_storages" "all" {
}

output "resp" {
  value = data.vmds_object_storages.all
}

