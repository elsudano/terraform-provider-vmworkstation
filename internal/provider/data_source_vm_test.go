// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVMDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccVMDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.vmworkstation_virtual_machine.parentvm",
						tfjsonpath.New("id"),
						knownvalue.StringExact("545OMDAL1R520604HKNKA6TTK6TBNOHK"),
					),
				},
			},
		},
	})
}

const testAccVMDataSourceConfig = `
provider "vmworkstation" {
  endpoint = "https://localhost:8697/api"
  username = "Admin"
  password = "Adm1n#01"
  https    = "true"
  debug    = "NONE"
}

data "vmworkstation_virtual_machine" "parentvm" {
  denomination = "parentvm"
}
`
