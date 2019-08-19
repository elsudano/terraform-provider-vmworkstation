package vmworkstation

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMWSFolder() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceVMWSFolderCreate,
		Read:   resourceVMWSFolderRead,
		Update: resourceVMWSFolderUpdate,
		Delete: resourceVMWSFolderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVMWSFolderImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
			},
		},
	}
}

func resourceVMWSFolderCreate(d *schema.ResourceData, m interface{}) error {
	return resourceVMWSFolderCreate(d, m)
}

func resourceVMWSFolderRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVMWSFolderUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVMWSFolderUpdate(d, m)
}

func resourceVMWSFolderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVMWSFolderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
