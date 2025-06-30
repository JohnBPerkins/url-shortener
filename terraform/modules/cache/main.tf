resource "aws_elasticache_subnet_group" "this" {
  name      = "${var.cluster_id}-subnet-group"
  subnet_ids = var.subnet_ids
}
resource "aws_elasticache_cluster" "this" {
  cluster_id           = var.cluster_id
  engine               = "redis"
  node_type            = var.node_type
  num_cache_nodes      = 1
  subnet_group_name    = aws_elasticache_subnet_group.this.name
}