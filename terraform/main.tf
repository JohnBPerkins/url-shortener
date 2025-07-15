# terraform {
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = "~> 4.0"
#     }
#   }
# }

provider "aws" {
  region = var.aws_region
}

########################
# Variables
########################

variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Deployment environment (e.g. dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnets" {
  description = "List of public subnet CIDRs"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnets" {
  description = "List of private subnet CIDRs"
  type        = list(string)
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "db_username" {
  description = "Master username for RDS"
  type        = string
  default     = "shortener"
}

variable "db_password" {
  description = "Master password for RDS"
  type        = string
}

variable "db_name" {
  description = "Initial database name"
  type        = string
  default     = "shortener"
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "cache_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "container_image_tag" {
  description = "Docker image tag for the service"
  type        = string
  default     = "latest"
}

variable "ecr_repository_name" {
  description = "ECR repository name"
  type        = string
  default     = "url-shortener"
}

########################
# Data Sources
########################

data "aws_caller_identity" "current" {}

data "aws_availability_zones" "available" {}

########################
# Networking
########################

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags                 = { Name = "${var.environment}-vpc" }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags   = { Name = "${var.environment}-igw" }
}

# Public subnets
resource "aws_subnet" "public" {
  for_each                = toset(var.public_subnets)
  vpc_id                  = aws_vpc.main.id
  cidr_block              = each.value
  map_public_ip_on_launch = true
  availability_zone       = element(data.aws_availability_zones.available.names, index(var.public_subnets, each.value))
  tags                    = { Name = "${var.environment}-public-${each.key}" }
}

# Private subnets
resource "aws_subnet" "private" {
  for_each                = toset(var.private_subnets)
  vpc_id                  = aws_vpc.main.id
  cidr_block              = each.value
  map_public_ip_on_launch = false
  availability_zone       = element(data.aws_availability_zones.available.names, index(var.private_subnets, each.value))
  tags                    = { Name = "${var.environment}-private-${each.key}" }
}

# NAT Gateway for private subnets
resource "aws_eip" "nat" {
  domain = "vpc"
  tags   = { Name = "${var.environment}-nat-eip" }
}

resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = values(aws_subnet.public)[0].id
  tags          = { Name = "${var.environment}-nat-gw" }
  depends_on    = [aws_internet_gateway.main]
}

# Route tables
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  tags   = { Name = "${var.environment}-public-rt" }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  tags   = { Name = "${var.environment}-private-rt" }

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main.id
  }
}

# Route table associations
resource "aws_route_table_association" "public" {
  for_each       = aws_subnet.public
  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  for_each       = aws_subnet.private
  subnet_id      = each.value.id
  route_table_id = aws_route_table.private.id
}

# Security groups
resource "aws_security_group" "app" {
  name        = "${var.environment}-app-sg"
  description = "Allow traffic from ALB and all egress"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "HTTP from ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    description     = "gRPC from ALB"
    from_port       = 50051
    to_port         = 50051
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = { Name = "${var.environment}-app-sg" }
}

resource "aws_security_group" "db" {
  name        = "${var.environment}-db-sg"
  description = "Allow PostgreSQL access from app"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "PostgreSQL"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  ingress {
    description = "PostgreSQL from developer IPv4"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.developer_ipv4]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = { Name = "${var.environment}-db-sg" }
}

resource "aws_security_group" "redis" {
  name        = "${var.environment}-redis-sg"
  description = "Allow Redis access from app"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Redis"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = { Name = "${var.environment}-redis-sg" }
}

########################
# Database (RDS Postgres)
########################

variable "developer_ipv4" {
  description = "Developer IPv4 address for database access"
  type        = string
  default     = "68.226.21.62/32"
}

resource "aws_db_subnet_group" "postgres_public" {
  name       = "${var.environment}-db-subnets-public"
  subnet_ids = [for s in aws_subnet.public : s.id]
  tags       = { Name = "${var.environment}-db-subnets-public" }
}

resource "aws_db_instance" "postgres" {
  identifier             = "${var.environment}-postgres"
  engine                 = "postgres"
  instance_class         = var.db_instance_class
  allocated_storage      = 20
  storage_type           = "gp2"
  db_name                = var.db_name
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.postgres_public.name
  vpc_security_group_ids = [aws_security_group.db.id]
  skip_final_snapshot    = true
  publicly_accessible    = true
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  tags = { Name = "${var.environment}-postgres" }
}

########################
# Cache (ElastiCache Redis)
########################

resource "aws_elasticache_subnet_group" "redis" {
  name       = "${var.environment}-redis-subnets"
  subnet_ids = [for s in aws_subnet.private : s.id]
}

resource "aws_elasticache_parameter_group" "redis_custom" {
  name        = "${var.environment}-redis-pg"
  family      = "redis7"
  description = "Redis parameter group with allkeys-lru"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.environment}-redis"
  engine               = "redis"
  engine_version       = "7.0"
  node_type            = var.cache_node_type
  num_cache_nodes      = 1
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.redis.id]
  port                 = 6379
  parameter_group_name = aws_elasticache_parameter_group.redis_custom.name
  
  tags = { Name = "${var.environment}-redis" }
}

########################
# Logging
########################

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/url-shortener-${var.environment}"
  retention_in_days = 7
  
  tags = { Name = "${var.environment}-ecs-logs" }
}

########################
# Application Load Balancer
########################

resource "aws_lb" "main" {
  name               = "${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = [for s in aws_subnet.public : s.id]

  enable_deletion_protection = false

  tags = { Name = "${var.environment}-alb" }
}

resource "aws_lb" "grpc" {
  name               = "${var.environment}-grpc-nlb"
  internal           = false
  load_balancer_type = "network"
  subnets            = [for s in aws_subnet.public : s.id]

  enable_deletion_protection = false

  tags = { Name = "${var.environment}-grpc-nlb" }
}

resource "aws_security_group" "alb" {
  name        = "${var.environment}-alb-sg"
  description = "Allow HTTP and HTTPS traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.environment}-alb-sg" }
}

resource "aws_lb_target_group" "grpc_nlb" {
  name     = "${var.environment}-grpc-nlb-tg"
  port     = 50051
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 10
    unhealthy_threshold = 2
  }

  tags = { Name = "${var.environment}-grpc-nlb-tg" }
}

resource "aws_lb_listener" "grpc_nlb" {
  load_balancer_arn = aws_lb.grpc.arn
  port              = "50051"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.grpc_nlb.arn
  }
}

# Target Group for Web Frontend
resource "aws_lb_target_group" "web" {
  name     = "${var.environment}-web-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }

  tags = { Name = "${var.environment}-web-tg" }
}

# ALB Listener for HTTP (Web Frontend)
resource "aws_lb_listener" "web" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.web.arn
  }
}

########################
# ECS (Fargate)
########################

# IAM role for task execution
data "aws_iam_policy_document" "task_exec" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "task_exec" {
  name               = "${var.environment}-ecs-exec"
  assume_role_policy = data.aws_iam_policy_document.task_exec.json
  
  tags = { Name = "${var.environment}-ecs-exec-role" }
}

resource "aws_iam_role_policy_attachment" "exec_policy" {
  role       = aws_iam_role.task_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.environment}-ecs-cluster"
  
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
  
  tags = { Name = "${var.environment}-ecs-cluster" }
}

# Task Definition
resource "aws_ecs_task_definition" "app" {
  family                   = "${var.environment}-shortener"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.task_exec.arn

  container_definitions = jsonencode([
    {
      name      = "shortener"
      image     = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.ecr_repository_name}:${var.container_image_tag}"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        },
        {
          containerPort = 50051
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "DATABASE_DSN"
          value = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.postgres.endpoint}/${var.db_name}?sslmode=require"
        },
        {
          name  = "REDIS_ENDPOINT"
          value = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:${aws_elasticache_cluster.redis.cache_nodes[0].port}"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.ecs.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
  
  tags = { Name = "${var.environment}-shortener-task" }
}

# ECS Service
resource "aws_ecs_service" "app" {
  name            = "${var.environment}-service"
  cluster         = aws_ecs_cluster.main.id
  launch_type     = "FARGATE"
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1
  force_new_deployment = true

  network_configuration {
    subnets          = [for s in aws_subnet.public : s.id]
    security_groups  = [aws_security_group.app.id]
    assign_public_ip = true
  }

  # Web frontend via ALB
  load_balancer {
    target_group_arn = aws_lb_target_group.web.arn
    container_name   = "shortener"
    container_port   = 8080
  }

  # gRPC via NLB (no SSL required)
  load_balancer {
    target_group_arn = aws_lb_target_group.grpc_nlb.arn
    container_name   = "shortener"
    container_port   = 50051
  }

  depends_on = [
    aws_iam_role_policy_attachment.exec_policy,
    aws_db_instance.postgres,
    aws_elasticache_cluster.redis,
    aws_lb_listener.web,
    aws_lb_listener.grpc_nlb
  ]
  
  tags = { Name = "${var.environment}-ecs-service" }
}

########################
# Outputs
########################

output "grpc_endpoint" {
  description = "gRPC endpoint via Network Load Balancer"
  value       = "${aws_lb.grpc.dns_name}:50051"
}

output "grpc_nlb_dns" {
  description = "DNS name of the gRPC Network Load Balancer"
  value       = aws_lb.grpc.dns_name
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = [for s in aws_subnet.public : s.id]
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = [for s in aws_subnet.private : s.id]
}

output "db_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.postgres.endpoint
}

output "redis_endpoint" {
  description = "Redis cluster endpoint"
  value       = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:${aws_elasticache_cluster.redis.cache_nodes[0].port}"
}

output "ecs_cluster_arn" {
  description = "ARN of the ECS cluster"
  value       = aws_ecs_cluster.main.arn
}

output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "web_frontend_url" {
  description = "URL to access the web frontend"
  value       = "http://${aws_lb.main.dns_name}"
}

output "ecr_image_uri" {
  description = "Full ECR image URI being used"
  value       = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.ecr_repository_name}:${var.container_image_tag}"
}