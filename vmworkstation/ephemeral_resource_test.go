package vmworkstation

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccVMEphemeralResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// Ephemeral resources are only available in 1.10 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			{
				Config: testAccVMEphemeralResourceConfig("example"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("token-123"),
					),
				},
			},
		},
	})
}

func testAccVMEphemeralResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
ephemeral "scaffolding_example" "test" {
  configurable_attribute = %[1]q
}

provider "echo" {
  data = ephemeral.scaffolding_example.test
}

resource "echo" "test" {}
`, configurableAttribute)
}
