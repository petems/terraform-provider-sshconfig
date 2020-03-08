package sshconfig

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider declaration of
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"sshconfig_host": dataSourceHost(),
		},

		ResourcesMap: map[string]*schema.Resource{},
	}
}
