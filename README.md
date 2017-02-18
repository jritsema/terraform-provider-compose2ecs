### compose2ecs

A tool that transforms a `docker-compose.yml` into a `task-definition.json`.

Uses docker labels for properties that don't translate to docker compose.

input 

```
version: "2"
services:  
  web:
    container_name: web
    image: 618440173123.dkr.ecr.us-east-1.amazonaws.com/web:$VERSION
    ports:
      - 80:80
    labels: 
      compose2ecs.hostPort: "0"
      compose2ecs.memoryReservation: "1000"
```

output

```
[
  {
    "Name": "web",    
    "Image": "618440173123.dkr.ecr.us-east-1.amazonaws.com/web:1.0",
    "MemoryReservation": 1000,
    "PortMappings": [
      {
        "ContainerPort": 80,
        "HostPort": 0
      }
    ]
  }
]
```

-----

TODO: convert to a terraform plugin that has a `compose2ecs` datasource, so that this functionality can be done 100% in terraform.

```terraform
data "compose2ecs" "compose" {
  compose_file = "${var.compose_file}"
}

resource "aws_ecs_task_definition" "app" {
  family                = "${var.app}"
  container_definitions = "${data.compose2ecs.compose.container_definitions}"
}
```