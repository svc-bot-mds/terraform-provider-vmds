data "vmds_clusters" "all_rmq" {
  service_type = "RABBITMQ"
}
data "vmds_clusters" "all_postgres" {
  service_type = "POSTGRES"
}
data "vmds_clusters" "all_redis" {
  service_type = "REDIS"
}
data "vmds_clusters" "all_mysql" {
  service_type = "MYSQL"
}
