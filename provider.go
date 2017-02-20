package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider returns a terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"compose2ecs": dataSourceCompose2Ecs(),
		},
	}
}
