package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
)

func transformComposeFile(composeFile string) (ecs.TaskDefinition, error) {

	//change compose version to 2 so that libcompose can parse it
	dat, err := ioutil.ReadFile(composeFile)
	if err != nil {
		log.Fatal(err)
	}

	temp := strings.Replace(string(dat), "version: '2.1'", "version: '2'", -1)
	temp = strings.Replace(temp, "version: \"2.1\"", "version: \"2\"", -1)

	ctx := ctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{
				[]byte(temp),
			},
		},
	}

	task := ecs.TaskDefinition{
		ContainerDefinitions: []*ecs.ContainerDefinition{},
		Volumes:              []*ecs.Volume{},
	}

	dockerComposeProject, err := docker.NewProject(&ctx, nil)
	if err != nil {
		return task, err
	}

	proj := dockerComposeProject.(*project.Project)
	for _, serviceName := range proj.ServiceConfigs.Keys() {

		config, success := dockerComposeProject.GetServiceConfig(serviceName)
		if !success {
			log.Fatal("error getting service config")
		}

		var def = ecs.ContainerDefinition{}

		if config.ContainerName != "" {
			def.Name = aws.String(config.ContainerName)
		}

		if config.Image != "" {
			def.Image = aws.String(config.Image)
		}

		if config.Hostname != "" {
			def.Hostname = aws.String(config.Hostname)
		}

		if config.WorkingDir != "" {
			def.WorkingDirectory = aws.String(config.WorkingDir)
		}

		if config.Privileged {
			def.Privileged = aws.Bool(config.Privileged)
		}

		if slice := config.DNS; len(slice) > 0 {
			for _, dns := range slice {
				def.DnsServers = append(def.DnsServers, aws.String(dns))
			}
		}

		if slice := config.DNSSearch; len(slice) > 0 {
			for _, item := range slice {
				def.DnsSearchDomains = append(def.DnsSearchDomains, aws.String(item))
			}
		}

		if cmds := config.Command; len(cmds) > 0 {
			def.Command = []*string{}
			for _, command := range cmds {
				def.Command = append(def.Command, aws.String(command))
			}
		}

		if cmds := config.Entrypoint; len(cmds) > 0 {
			def.EntryPoint = []*string{}
			for _, command := range cmds {
				def.EntryPoint = append(def.EntryPoint, aws.String(command))
			}
		}

		if slice := config.Environment; len(slice) > 0 {
			def.Environment = []*ecs.KeyValuePair{}
			for _, val := range slice {
				parts := strings.SplitN(val, "=", 2)
				def.Environment = append(def.Environment, &ecs.KeyValuePair{
					Name:  aws.String(parts[0]),
					Value: aws.String(parts[1]),
				})
			}
		}

		if ports := config.Ports; len(ports) > 0 {
			def.PortMappings = []*ecs.PortMapping{}
			for _, val := range ports {
				parts := strings.Split(val, ":")
				mapping := &ecs.PortMapping{}

				// TODO: support host to map to
				if len(parts) > 0 {
					portInt, err := strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						return task, err
					}
					mapping.ContainerPort = aws.Int64(portInt)
				}

				if len(parts) > 1 {
					hostParts := strings.Split(parts[1], "/")
					portInt, err := strconv.ParseInt(hostParts[0], 10, 64)
					if err != nil {
						return task, err
					}
					mapping.HostPort = aws.Int64(portInt)

					hostPortOverride := config.Labels["compose2ecs.hostPort"]
					if len(hostPortOverride) > 0 {
						temp, err := strconv.Atoi(hostPortOverride)
						if err != nil {
							return task, err
						}
						mapping.HostPort = aws.Int64(int64(temp))
					}

					// handle the protocol at the end of the mapping
					if len(hostParts) > 1 {
						mapping.Protocol = aws.String(hostParts[1])
					}
				}

				if len(parts) == 0 || len(parts) > 2 {
					return task, errors.New("Unsupported port mapping " + val)
				}

				def.PortMappings = append(def.PortMappings, mapping)
			}
		}

		if links := config.Links; len(links) > 0 {
			def.Links = []*string{}
			for _, link := range links {
				def.Links = append(def.Links, aws.String(link))
			}
		}

		if volsFrom := config.VolumesFrom; len(volsFrom) > 0 {
			def.VolumesFrom = []*ecs.VolumeFrom{}
			for _, container := range volsFrom {
				def.VolumesFrom = append(def.VolumesFrom, &ecs.VolumeFrom{
					SourceContainer: aws.String(container),
				})
			}
		}

		memoryReservationOverride := config.Labels["compose2ecs.memoryReservation"]
		if len(memoryReservationOverride) > 0 {
			temp, err := strconv.Atoi(memoryReservationOverride)
			if err != nil {
				return task, err
			}
			def.MemoryReservation = aws.Int64(int64(temp))
		}

		task.ContainerDefinitions = append(task.ContainerDefinitions, &def)
	}

	return task, nil
}
