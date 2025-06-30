variable "db_username" { type = string }
variable "db_password" { type = string }
variable "db_name"     { type = string }
variable "engine"      { type = string }
variable "instance_class" { type = string }
variable "subnet_ids"  { type = list(string) }
variable "security_group_ids" { type = list(string) }
variable "environment" { type = string }