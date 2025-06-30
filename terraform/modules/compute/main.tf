resource "aws_ecs_cluster" "this" {
  name = var.cluster_name
}
resource "aws_iam_role" "task_exec" {
  name = "${var.cluster_name}-exec-role"
  assume_role_policy = data.aws_iam_policy_document.task_exec.json
}
# ... define policy attachment and task definition, service ...