variable "cluster_name" { type = string }
variable "subnet_ids"    { type = list(string) }
variable "security_group_ids" { type = list(string) }
variable "container_image" { type = string }