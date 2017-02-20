### terraform-provider-compose2ecs

A terraform plugin containing a datasource that can transform a docker compose file into an ecs task defnition.

#### usage

```terraform
data "compose2ecs" "compose" {}

resource "aws_ecs_task_definition" "app" {
  family                = "${var.app}"
  container_definitions = "${data.compose2ecs.compose.container_definitions}"
}
```

*note that you can specify `compose_file` if you want to override the default compose file name (`docker-compose.yml`), for example...

```
data "compose2ecs" "compose" {
  compose_file = "my-compose.yml"
}
```

where `docker-compose.yml` might look like...

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

and the outputted container_definitions would be...

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
