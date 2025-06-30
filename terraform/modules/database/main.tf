resource "aws_db_subnet_group" "this" {
  name       = "${var.db_name}-subnet-group"
  subnet_ids = var.subnet_ids
}
resource "aws_db_instance" "this" {
  identifier         = "${var.db_name}-${var.environment}"
  engine             = var.engine
  instance_class     = var.instance_class
  username           = var.db_username
  password           = var.db_password
  db_subnet_group_name = aws_db_subnet_group.this.name
  vpc_security_group_ids = var.security_group_ids
  allocated_storage  = 20
  skip_final_snapshot = true
}