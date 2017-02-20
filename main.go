package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

func main2() {

	taskDefinition, err := transformComposeFile("docker-compose.yml")
	if err != nil {
		log.Fatal(err)
	}

	//serialize object to json
	byteArray, err := json.MarshalIndent(taskDefinition.ContainerDefinitions, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding to JSON: %s", err)
	}

	//write output to stdout and file
	fmt.Println(string(byteArray))
	err = ioutil.WriteFile("task-definition.json", byteArray, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
