package vmworkstation

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"example": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VMWS_USER"); v == "" {
		t.Fatal("VMWS_USER must be set for acceptance tests")
	}
	if v := os.Getenv("VMWS_PASSWORD"); v == "" {
		t.Fatal("VMWS_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("VMWS_URL"); v == "" {
		t.Fatal("VMWS_URL must be set for acceptance tests")
	}
}
