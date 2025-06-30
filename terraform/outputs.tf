output "vpc_id" {
  value = module.network.vpc_id
}
output "db_endpoint" {
  value = module.database.db_endpoint
}
output "redis_endpoint" {
  value = module.cache.redis_endpoint
}
output "ecs_cluster_arn" {
  value = module.compute.cluster_arn
}