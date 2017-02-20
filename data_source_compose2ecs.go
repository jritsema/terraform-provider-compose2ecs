package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCompose2Ecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCompose2EcsRead,

		Schema: map[string]*schema.Schema{
			"compose_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "docker-compose.yml",
			},
			"services": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"container_definitions": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomID() string {
	b := make([]rune, 30)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func dataSourceCompose2EcsRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(generateRandomID())

	composeFile := d.Get("compose_file").(string)

	var services []string
	set := d.Get("services").(*schema.Set)
	for _, e := range set.List() {
		s := e.(string)
		set.Remove(s)
		services = append(services, strings.ToLower(s))
	}

	//transform docker compose file into an ecs task definition
	taskDefinition, err := transformComposeFile(composeFile, services)
	if err != nil {
		log.Fatal(err)
	}

	//serialize object to json
	byteArray, err := json.MarshalIndent(taskDefinition.ContainerDefinitions, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding to JSON: %s", err)
	}

	//output task definition json
	d.Set("container_definitions", string(byteArray))

	return nil
}
