package vmworkstation

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMWSFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMWSFolderCreate,
		Read:   resourceVMWSFolderRead,
		Update: resourceVMWSFolderUpdate,
		Delete: resourceVMWSFolderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVMWSFolderImport,
		},

		SchemaVersion: 1,
		Schema:        map[string]*schema.Schema{},
	}
}

func resourceVMWSFolderCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceVMWSFolderCreate(d, meta)
}

func resourceVMWSFolderRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVMWSFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceVMWSFolderUpdate(d, meta)
}

func resourceVMWSFolderDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVMWSFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
