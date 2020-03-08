package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/petems/terraform-provider-sshconfig/sshconfig"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sshconfig.Provider})
}
