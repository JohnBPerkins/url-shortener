variable "aws_region" {
  type    = string
  default = "us-east-1"
}
variable "environment" {
  type    = string
  default = "dev"
}
variable "state_bucket" {
  type = string
}
variable "lock_table" {
  type = string
}