resource "aws_vpc" "this" {
  cidr_block = var.vpc_cidr
  tags = { Name = "${var.environment}-vpc" }
}
resource "aws_subnet" "public" {
  for_each = toset(var.public_subnets)
  vpc_id            = aws_vpc.this.id
  cidr_block        = each.value
  map_public_ip_on_launch = true
  tags = { Name = "${var.environment}-public-${each.key}" }
}
resource "aws_subnet" "private" {
  for_each = toset(var.private_subnets)
  vpc_id     = aws_vpc.this.id
  cidr_block = each.value
  tags = { Name = "${var.environment}-private-${each.key}" }
}
resource "aws_security_group" "app" {
  name        = "${var.environment}-app-sg"
  description = "Allow HTTP and gRPC"
  vpc_id      = aws_vpc.this.id
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 50051
    to_port     = 50051
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress { 
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}