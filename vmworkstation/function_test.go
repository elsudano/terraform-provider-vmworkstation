package vmworkstation

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestVMFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				provider "vmworkstation" {
					endpoint = "https://192.168.1.155:8697/api"
					username = "Admin"
					password = "Adm1n#01"
					https    = "true"
					debug    = "NONE"
				}

				output "test" {
					value = provider::vmworkstation::vm1("testvalue")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact("testvalue"),
					),
				},
			},
		},
	})
}

func TestVMFunction_Null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				provider "vmworkstation" {
					endpoint = "https://192.168.1.155:8697/api"
					username = "Admin"
					password = "Adm1n#01"
					https    = "true"
					debug    = "NONE"
				}

				output "test" {
					value = provider::vmworkstation::vm1(null)
				}
				`,
				// The parameter does not enable AllowNullValue
				ExpectError: regexp.MustCompile(`argument must not be null`),
			},
		},
	})
}

func TestVMFunction_Unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				provider "vmworkstation" {
					endpoint = "https://192.168.1.155:8697/api"
					username = "Admin"
					password = "Adm1n#01"
					https    = "true"
					debug    = "NONE"
				}

				resource "terraform_data" "test" {
					input = "testvalue"
				}
				
				output "test" {
					value = provider::vmworkstation::vm1(terraform_data.test.output)
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact("testvalue"),
					),
				},
			},
		},
	})
}
