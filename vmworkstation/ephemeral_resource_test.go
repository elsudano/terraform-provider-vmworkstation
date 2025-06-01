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
ephemeral "vmworkstation_ephemeral" "vm1" {
  configurable_attribute = %[1]q
}

provider "vmworkstation" {
  endpoint = "https://localhost:8697/api"
  username = "Admin"
  password = "Adm1n#01"
  https    = "true"
  debug    = "NONE"
}

provider "echo" {
  data = ephemeral.vmworkstation_ephemeral.vm1
}

resource "echo" "test" {}
`, configurableAttribute)
}
