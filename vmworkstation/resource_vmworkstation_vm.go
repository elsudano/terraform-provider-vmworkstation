package vmworkstation

import (
	"log"
	"strings"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMWSVm() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceVMWSVmCreate,
		Read:   resourceVMWSVmRead,
		Update: resourceVMWSVmUpdate,
		Delete: resourceVMWSVmDelete,
		Exists: resourceVMWSVmExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sourceid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the VM that to use for clone at the new",
			},
			"denomination": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Little bit description of the instance",
			},
			// "image": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	Description: "String with ID for the image that create the VM",
			// },
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

func resourceVMWSVmCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	sourceid := d.Get("sourceid").(string)
	denomination := d.Get("denomination").(string)
	description := d.Get("description").(string)
	// image := d.Get("image").(string)
	processors := d.Get("processors").(int)
	memory := d.Get("memory").(int)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmCreate Ob:\n")
	log.Printf("\t sourceid: %#v\n", sourceid)
	log.Printf("\t denomination: %#v\n", denomination)
	log.Printf("\t description: %#v\n", description)
	// log.Printf("\t image: %#v\n", image)
	log.Printf("\t processors: %#v\n", processors)
	log.Printf("\t memory: %#v\n", memory)
	VM, err := apiClient.CreateVM(sourceid, denomination)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(VM.IdVM)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmCreate Obj:ID of new VM %#v\n", VM.IdVM)
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	MyVM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(MyVM.IdVM)
	d.Set("denomination", MyVM.Denomination)
	d.Set("description", MyVM.Description)
	d.Set("image", MyVM.Image)
	d.Set("processors", MyVM.CPU.Processors)
	d.Set("memory", MyVM.Memory)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmRead Obj:One VM %#v\n", MyVM)
	return nil
}

func resourceVMWSVmUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true) // this is such as to a semaphore, it's a switch to change a state of blocked
	if d.HasChange("denomination") {
		d.SetId("denomination")
		d.SetPartial("denomination")
	}
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmUpdate Obj:DataScheme %#v\n", d)
	d.Partial(false) // this is such as to a semaphore, it's a switch to change a state of unblocked
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	err := apiClient.DeleteVM(d.Id())
	if err != nil {
		log.Printf("[VMWS][ERROR] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmDelete Obj:%#v\n", err)
		return nil
	}
	d.SetId("")
	return nil
}

func resourceVMWSVmExists(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*wsapiclient.Client)
	MyVM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	log.Printf("[VMWS] Fi: resource_vmworkstation_vms.go Fu: resourceVMWSVmExists Obj:APIClient %#v\n", apiClient)
	if MyVM == nil {
		return false, nil
	}
	return true, nil
}

func resourceVMWSVmImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
