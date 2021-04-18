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
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path absolute of the VM machine",
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
	path := d.Get("path").(string)
	// image := d.Get("image").(string)
	processors := d.Get("processors").(int)
	memory := d.Get("memory").(int)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob:\n")
	log.Printf("\t sourceid: %#v\n", sourceid)
	log.Printf("\t denomination: %#v\n", denomination)
	log.Printf("\t description: %#v\n", description)
	// log.Printf("\t image: %#v\n", image)
	log.Printf("\t path: %#v\n", path)
	log.Printf("\t processors: %#v\n", processors)
	log.Printf("\t memory: %#v\n", memory)
	VM, err := apiClient.CreateVM(sourceid, denomination, description)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob: %#v\n", VM.IdVM)
	if err != nil {
		d.SetId("")
		return nil
	}
	VM, err = apiClient.UpdateVM(VM.IdVM, denomination, description, processors, memory)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob: %#v\n", VM.IdVM)
	if err != nil {
		d.SetId("")
		return nil
	}
	VM, err = apiClient.RegisterVM(denomination, path)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob: %#v\n", VM.IdVM)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(VM.IdVM)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Obj:ID of new VM %#v\n", VM.IdVM)
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	VM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(VM.IdVM)
	d.Set("denomination", VM.Denomination)
	d.Set("description", VM.Description)
	// d.Set("image", VM.Image)
	d.Set("processors", VM.CPU.Processors)
	d.Set("memory", VM.Memory)
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmRead Obj:One VM %#v\n", VM)
	return nil
}

func resourceVMWSVmUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	VM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: VM to Update %#v\n", VM)
	d.Partial(true) // this is such as to a semaphore, it's a switch to change a state of blocked
	if d.HasChange("denomination") {
		DenominationOldState, DenominationNewState := d.GetChange("denomination")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Denomination %#v\n", DenominationOldState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Denomination %#v\n", DenominationNewState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: denomination field in VM before %#v\n", VM.Denomination)
		if VM, err = apiClient.UpdateVM(d.Id(), d.Get("denomination").(string), VM.Description, VM.CPU.Processors, VM.Memory); err != nil {
			return nil
		}
		d.SetPartial("denomination")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: denomination field in VM after %#v\n", VM.Denomination)
	}
	if d.HasChange("description") {
		DescriptionOldState, DescriptionNewState := d.GetChange("description")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Description %#v\n", DescriptionOldState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Description %#v\n", DescriptionNewState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Description field in VM before %#v\n", VM.Description)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, d.Get("description").(string), VM.CPU.Processors, VM.Memory); err != nil {
			return nil
		}
		d.SetPartial("denomination")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Description field in VM after %#v\n", VM.Description)
	}
	if d.HasChange("processors") {
		ProcessorsOldState, ProcessorsNewState := d.GetChange("processors")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Processors %#v\n", ProcessorsOldState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Processors %#v\n", ProcessorsNewState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Processors field in VM before %#v\n", VM.CPU.Processors)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, VM.Description, d.Get("processors").(int), VM.Memory); err != nil {
			return nil
		}
		d.SetPartial("processors")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Processors field in VM after %#v\n", VM.CPU.Processors)
	}
	if d.HasChange("memory") {
		MemoryOldState, MemoryNewState := d.GetChange("memory")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Memory %#v\n", MemoryOldState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Memory %#v\n", MemoryNewState)
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Memory field in VM before %#v\n", VM.Memory)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, VM.Description, VM.CPU.Processors, d.Get("memory").(int)); err != nil {
			return nil
		}
		d.SetPartial("memory")
		log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Memory field in VM after %#v\n", VM.Memory)
	}
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj:DataScheme %#v\n", d)
	d.Partial(false) // this is such as to a semaphore, it's a switch to change a state of unblocked
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	err := apiClient.DeleteVM(d.Id())
	if err != nil {
		log.Printf("[VMWS][ERROR] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmDelete Obj:%#v\n", err)
		return nil
	}
	d.SetId("")
	return nil
}

func resourceVMWSVmExists(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*wsapiclient.Client)
	VM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	log.Printf("[VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmExists Obj:APIClient %#v\n", apiClient)
	if VM == nil {
		return false, nil
	}
	return true, nil
}

func resourceVMWSVmImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
