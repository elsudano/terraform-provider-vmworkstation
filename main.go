package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/elsudano/terraform-provider-vmworkstation/vmworkstation"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}