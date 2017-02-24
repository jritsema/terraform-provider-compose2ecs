### terraform-provider-compose2ecs

A [terraform](terraform.io) plugin containing a datasource that can transform a [docker compose file](https://docs.docker.com/compose/compose-file/) into an [ecs task definition](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_defintions.html).

![CircleCI](https://circleci.com/gh/jritsema/terraform-provider-compose2ecs.svg?style=shield&circle-token=:circle-token)

#### usage

Download and install the [plugin](https://github.com/jritsema/terraform-provider-compose2ecs/releases)

```
$ wget -O /usr/local/bin/terraform-provider-compose2ecs https://github.com/jritsema/terraform-provider-compose2ecs/releases/download/v0.1.0-1-g7bcb595/ncd_darwin_amd64 && chmod +x /usr/local/bin/terraform-provider-compose2ecs
```

```terraform
data "compose2ecs" "compose" {}

resource "aws_ecs_task_definition" "app" {
  family                = "${var.app}"
  container_definitions = "${data.compose2ecs.compose.container_definitions}"
}
```

You can optionally specify `compose_file` if you want to override the default compose file name (defaults to `docker-compose.yml`).  You can also optionally specify which subset of services from the compose file you want to include (defaults to all).

```terraform
data "compose2ecs" "compose" {
  compose_file = "my-compose.yml"
  services     = ["web", "worker"]
}
```

where `docker-compose.yml` might look like...

```yaml
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

and the outputted `container_definitions` would be...

```json
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
