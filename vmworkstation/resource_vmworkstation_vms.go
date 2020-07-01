package vmworkstation

import (
	"log"
	"strings"

	wsapiclient "github.com/elsudano/vmware-workstation-api-client/wsapiclient"
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
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "String with ID for the image that create the VM",
			},
			"processors": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of processors that will have the VM",
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "How much memory will have the VM",
			},
		},
	}
}

func resourceVMWSVmsCreate(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*wsapiclient.Client)

	id_name := d.Get("name").(string)
	// apiClient.GetVM("id_name")
	d.SetId(id_name)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmsCreate Ob: %#v\n", id_name)
	return resourceVMWSVmsRead(d, m)
}

func resourceVMWSVmsRead(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*wsapiclient.Client)

	// MyVM, err := apiClient.GetVM(d.Id())

	// if err != nil {
	// 	d.SetId("")
	// 	return nil
	// }

	// d.Set("name", MyVM.IdVM)
	return nil
}

func resourceVMWSVmsUpdate(d *schema.ResourceData, m interface{}) error {
	// d.Partial(true)

	// if d.HasChange("name") {
	// 	d.SetId("name")
	// }
	// d.SetPartial("name")

	// d.Partial(false)
	return resourceVMWSVmsRead(d, m)
}

func resourceVMWSVmsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
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
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmsExists Ob: %#v\n", apiClient)
	return true, nil
}

func resourceVMWSVmsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
