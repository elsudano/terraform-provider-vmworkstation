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
				Config: testAccVMResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("example value when not configured"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("one"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "vmworkstation_resource.vm1",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccVMResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("example value when not configured"),
					),
					statecheck.ExpectKnownValue(
						"vmworkstation_resource.vm1",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("two"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccVMResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "vmworkstation_resource" "vm1" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
