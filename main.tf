provider "aws" {
  region = "us-east-1"
}

data "compose2ecs" "compose" {}

resource "aws_ecs_task_definition" "app" {
  family                = "compose2ecs-test"
  container_definitions = "${data.compose2ecs.compose.container_definitions}"
}
