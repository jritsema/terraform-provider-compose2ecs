package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

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
