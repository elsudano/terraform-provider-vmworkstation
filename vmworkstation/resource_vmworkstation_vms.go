package vmworkstation

import (
	"log"
	"strings"

	wsapiclient "github.com/elsudano/vmware-workstation-api-client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMWSVms() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceVMWSVmsCreate,
		Read:   resourceVMWSVmsRead,
		Update: resourceVMWSVmsUpdate,
		Delete: resourceVMWSVmsDelete,
		Exists: resourceVMWSVmsExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVMWSVmsCreate(d *schema.ResourceData, m interface{}) error {
	return resourceVMWSVmsCreate(d, m)
}

func resourceVMWSVmsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVMWSVmsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVMWSVmsUpdate(d, m)
}

func resourceVMWSVmsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVMWSVmsExists(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*wsapiclient.Client)

	VMId := d.Id()
	_, err := apiClient.GetVM(VMId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	if d.Get("debug") == true {
		log.Printf("[VMWS] Fu: resourceVMWSVmsExists Fi: resource_vmworkstation_vms.go Ob: %#v\n", apiClient)
	}
	return true, nil
}

func resourceVMWSVmsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
