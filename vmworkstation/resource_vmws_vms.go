package vmworkstation

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMWSVms() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMWSVmsCreate,
		Read:   resourceVMWSVmsRead,
		Update: resourceVMWSVmsUpdate,
		Delete: resourceVMWSVmsDelete,

		SchemaVersion: 1,
		Schema:        map[string]*schema.Schema{},
	}
}

func resourceVMWSVmsCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceVMWSVmsCreate(d, meta)
}

func resourceVMWSVmsRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVMWSVmsUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceVMWSVmsUpdate(d, meta)
}

func resourceVMWSVmsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
