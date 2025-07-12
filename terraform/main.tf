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

variable "container_image" {
  description = "Docker image URI for the service"
  type        = string
}

########################
# Networking
########################

# VPC
resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr
  tags = { Name = "${var.environment}-vpc" }
}

# Public subnets
resource "aws_subnet" "public" {
  for_each            = toset(var.public_subnets)
  vpc_id              = aws_vpc.main.id
  cidr_block          = each.value
  map_public_ip_on_launch = true
  availability_zone   = element(data.aws_availability_zones.available.names, index(var.public_subnets, each.value))
  tags = { Name = "${var.environment}-public-${each.key}" }
}

# Private subnets
resource "aws_subnet" "private" {
  for_each            = toset(var.private_subnets)
  vpc_id              = aws_vpc.main.id
  cidr_block          = each.value
  map_public_ip_on_launch = false
  availability_zone   = element(data.aws_availability_zones.available.names, index(var.private_subnets, each.value))
  tags = { Name = "${var.environment}-private-${each.key}" }
}

data "aws_availability_zones" "available" {}

# Security group
resource "aws_security_group" "app" {
  name        = "${var.environment}-app-sg"
  description = "Allow HTTP, gRPC, and all egress"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "gRPC"
    from_port   = 50051
    to_port     = 50051
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.environment}-app-sg" }
}

########################
# Database (RDS Postgres)
########################

resource "aws_db_subnet_group" "postgres" {
  name       = "${var.environment}-db-subnets"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_db_instance" "postgres" {
  identifier             = "${var.environment}-postgres"
  engine                 = "postgres"
  instance_class         = var.db_instance_class
  allocated_storage      = 20
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.postgres.name
  vpc_security_group_ids = [aws_security_group.app.id]
  skip_final_snapshot    = true
  publicly_accessible    = false
}

########################
# Cache (ElastiCache Redis)
########################

resource "aws_elasticache_subnet_group" "redis" {
  name       = "${var.environment}-redis-subnets"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.environment}-redis"
  engine               = "redis"
  node_type            = var.cache_node_type
  num_cache_nodes      = 1
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.app.id]
  port                 = 6379
}

########################
# Logging
########################

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/url-shortener"
  retention_in_days = 7
}

########################
# ECS (Fargate)
########################

# IAM role for task execution
data "aws_iam_policy_document" "task_exec" {
  statement {
    effect    = "Allow"
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
}

resource "aws_iam_role_policy_attachment" "exec_policy" {
  role       = aws_iam_role.task_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.environment}-ecs-cluster"
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
      image     = var.container_image
      essential = true
      portMappings = [
        {
          containerPort = 50051
          protocol      = "tcp"
        }
      ]
      environment = [
        { 
            name = "DATABASE_DSN"
            value = aws_db_instance.postgres.address 
        },
        { 
            name = "REDIS_ENDPOINT" 
            value = aws_elasticache_cluster.redis.primary_endpoint_address 
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
}

# ECS Service
resource "aws_ecs_service" "app" {
  name            = "${var.environment}-service"
  cluster         = aws_ecs_cluster.main.id
  launch_type     = "FARGATE"
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.app.id]
    assign_public_ip = false
  }

  depends_on = [aws_iam_role_policy_attachment.exec_policy]
}

########################
# Outputs
########################

output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_ids" {
  value = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  value = aws_subnet.private[*].id
}

output "db_endpoint" {
  value = aws_db_instance.postgres.address
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.redis.primary_endpoint_address
}

output "ecs_cluster_arn" {
  value = aws_ecs_cluster.main.arn
}
