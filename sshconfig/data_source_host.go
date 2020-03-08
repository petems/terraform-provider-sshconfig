package sshconfig

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gosshconfig "github.com/petems/go-sshconfig"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostRead,

		Schema: map[string]*schema.Schema{
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host information rendered as a multi-line string",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The host information as a map",
				Computed:    true,
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/etc/ssh/ssh_config",
				Description: "The path to the SSH config file you want to read",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host you want to lookup",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func checkHostExists(path string, host string) error {
	configFile, err := os.Open(path)

	if err != nil {
		fmt.Errorf("Error when opening file: %s", err)
	}

	config, err := gosshconfig.Parse(configFile)

	if err != nil {
		fmt.Errorf("Error when parsing ssh config: %s", err)
	}

	hostLookup := config.FindByHostname(host)

	if hostLookup == nil {
		fmt.Errorf("Could not find host %s in config", host)
	}

	return nil
}

func getHostFromSSHConfig(path string, host string) (string, error) {

	configFile, err := os.Open(path)

	if err != nil {
		return "", fmt.Errorf("Error when opening file: %s", err)
	}

	config, err := gosshconfig.Parse(configFile)

	if err != nil {
		return "", fmt.Errorf("Error when parsing ssh config: %s", err)
	}

	hostLookup := config.FindByHostname(host)

	if hostLookup == nil {
		return "", fmt.Errorf("Could not find host %s in config", host)
	}

	hostLookupString := hostLookup.String()

	return hostLookupString, nil
}

func getHashFromSSHConfig(path string, host string) (map[string]interface{}, error) {

	configFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Error when opening file: %s", err)
	}

	config, err := gosshconfig.Parse(configFile)

	if err != nil {
		return nil, fmt.Errorf("Error when parsing ssh config: %s", err)
	}

	hostLookup := config.FindByHostname(host)

	if hostLookup == nil {
		return nil, fmt.Errorf("Could not find host %s in config", host)
	}

	m := make(map[string]interface{})

	for _, param := range hostLookup.Params {
		log.Printf("[DEBUG] Param was %s=%s", param.Keyword, param.Args)
		m[param.Keyword] = strings.Join(param.Args, " ")
	}

	log.Printf("[DEBUG] Map was %s", m)

	return m, nil
}

func dataSourceHostRead(d *schema.ResourceData, meta interface{}) error {

	configPath := d.Get("path").(string)
	host := d.Get("host").(string)

	err := checkHostExists(configPath, host)

	if err != nil {
		return fmt.Errorf("Error reading ssh config: %d", err)
	} else {
		hostDeclaration, err := getHostFromSSHConfig(configPath, host)

		if err == nil {
			d.Set("rendered", string(hostDeclaration))
			d.SetId(time.Now().UTC().String())
		} else {
			return fmt.Errorf("Error reading ssh config: %d", err)
		}

		foo, err := getHashFromSSHConfig(configPath, host)

		log.Printf("[DEBUG] Map2 was %s", foo)

		if err != nil {
			return fmt.Errorf("Error reading ssh config: %d", err)
		}

		d.Set("host_map", foo)

	}

	return nil

}
