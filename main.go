// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package main

import (
	"context"
	"flag"
	"log"

	vmworkstation "github.com/elsudano/terraform-provider-vmworkstation/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := providerserver.ServeOpts{
		ProtocolVersion: 6,
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/elsudano/vmworkstation",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), vmworkstation.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
