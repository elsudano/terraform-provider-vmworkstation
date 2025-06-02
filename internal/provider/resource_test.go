// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVMResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVMResourceConfig(2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("sourceid"),
						knownvalue.StringExact("545OMDAL1R520604HKNKA6TTK6TBNOHK"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("processors"),
						knownvalue.Int32Exact(2),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("ip"),
						knownvalue.StringExact("0.0.0.0/0"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "vmworkstation_resource_vm.vm1",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"sourceid", "denomination", "description", "processors", "memory", "path", "state", "ip"},
			},
			// Update and Read testing
			{
				Config: testAccVMResourceConfig(4),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("sourceid"),
						knownvalue.StringExact("545OMDAL1R520604HKNKA6TTK6TBNOHK"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("processors"),
						knownvalue.Int32Exact(4),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource_vm.vm1",
						tfjsonpath.New("ip"),
						knownvalue.StringExact("0.0.0.0/0"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccVMResourceConfig(configurableAttribute int) string {
	return fmt.Sprintf(`
provider "vmworkstation" {
  endpoint = "https://192.168.1.155:8697/api"
  username = "Admin"
  password = "Adm1n#01"
  https    = "true"
  debug    = "NONE"
}

resource "vmworkstation_resource_vm" "vm1" {
  sourceid     = "545OMDAL1R520604HKNKA6TTK6TBNOHK"
  denomination = "go_tests_vm1"
  description  = "This VM is just a resource created by the GO tests."
  path         = "D:\\VirtualMachines\\vm1\\vm1.vmx"
  processors   = %[1]v
  memory       = 1024
  state        = "off"
}
`, configurableAttribute)
}
