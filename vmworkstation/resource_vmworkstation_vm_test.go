package vmworkstation

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("example_item.test_vm"),
					resource.TestCheckResourceAttr("example_item.test_vm", "name", "test"),
					resource.TestCheckResourceAttr("example_item.test_vm", "description", "hello"),
					// resource.TestCheckResourceAttr("example_item.test_vm", "tags.#", "2"),
					// resource.TestCheckResourceAttr("example_item.test_vm", "tags.1931743815", "tag1"),
					// resource.TestCheckResourceAttr("example_item.test_vm", "tags.1477001604", "tag2"),
				),
			},
		},
	})
}

func testAccCheckVMDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*wsapiclient.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "example_item" {
			continue
		}
		_, err := apiClient.CreateVM(rs.Primary.ID, rs.Primary.Attributes["name"], rs.Primary.Attributes["description"])
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckVMExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*wsapiclient.Client)
		_, err := apiClient.CreateVM(rs.Primary.ID, rs.Primary.Attributes["name"], rs.Primary.Attributes["description"])
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckVMBasic() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
  name        = "test"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "example_item" "test_update" {
  name        = "test_update"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "example_item" "test_update" {
  name        = "test_update"
  description = "updated description"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
`)
}

func testAccCheckItemMultiple() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
  name        = "test"
  description = "hello"
  
  tags = [
	"tag1",
    "tag2",
  ]
}
resource "example_item" "another_item" {
	name        = "another_test"
	description = "hello"
	
	tags = [
	  "tag1",
	  "tag2",
	]
  }
`)
}

func testAccCheckItemWhitespace() string {
	return fmt.Sprintf(`
resource "example_item" "test_item" {
	name        = "test with whitespace"
	description = "hello"
	tags = [
		"tag1",
		"tag2",
	]
}
`)
}
