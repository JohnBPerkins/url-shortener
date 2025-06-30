terraform {
  required_version = ">= 1.0.0"
  backend "s3" {
    bucket         = var.state_bucket
    key            = "url-shortener/${var.environment}.tfstate"
    region         = var.aws_region
    dynamodb_table = var.lock_table
    encrypt        = true
  }
}