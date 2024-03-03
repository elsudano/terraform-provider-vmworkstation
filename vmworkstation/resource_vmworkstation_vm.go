package vmworkstation

import (
	"log"
	"strings"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional:    true,
				Required:    false,
				Description: "The name of the resource",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
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
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Which will be the state of the VM when we will deploy it",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Which is the IP of the instance",
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
	state := d.Get("state").(string)
	ip := d.Get("ip").(string)
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob:\n")
	log.Printf("sourceid: %#v\n", sourceid)
	log.Printf("denomination: %#v\n", denomination)
	log.Printf("description: %#v\n", description)
	// log.Printf("image: %#v\n", image)
	log.Printf("path: %#v\n", path)
	log.Printf("processors: %#v\n", processors)
	log.Printf("memory: %#v\n", memory)
	log.Printf("state: %#v\n", state)
	log.Printf("IP: %#v\n", ip)
	VM, err := apiClient.CreateVM(sourceid, denomination, description, processors, memory)
	if err != nil {
		d.SetId("")
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Error Creating VM: %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob: %#v\n", VM.IdVM)
	VM, err = apiClient.PowerSwitch(VM.IdVM, state)
	if err != nil {
		d.SetId("")
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Error Powerized VM: %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Ob: After to Powerized %#v\n", VM.IdVM)
	VM, err = apiClient.RegisterVM(denomination, path)
	if err != nil {
		d.SetId("")
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Error Registering VM: %#v\n", err)
		return err
	}
	d.SetId(VM.IdVM)
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmCreate Obj:ID of new VM %#v\n", VM.IdVM)
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	VM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		d.SetId("")
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmRead Error Reading VM: %#v\n", err)
		return err
	}
	d.SetId(VM.IdVM)
	d.Set("denomination", VM.Denomination)
	d.Set("description", VM.Description)
	// d.Set("image", VM.Image)
	d.Set("processors", VM.CPU.Processors)
	d.Set("memory", VM.Memory)
	d.Set("state", VM.PowerStatus)
	d.Set("ip", VM.Ip)
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmRead Obj:One VM %#v\n", VM)
	return nil
}

func resourceVMWSVmUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*wsapiclient.Client)
	VM, err := apiClient.ReadVM(d.Id())
	if err != nil {
		d.SetId("")
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Reading VM: %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: VM to Update %#v\n", VM)
	d.Partial(true) // this is such as to a semaphore, it's a switch to change a state of blocked
	if d.HasChange("denomination") {
		DenominationOldState, DenominationNewState := d.GetChange("denomination")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Denomination %#v\n", DenominationOldState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Denomination %#v\n", DenominationNewState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: denomination field in VM before %#v\n", VM.Denomination)
		if VM, err = apiClient.UpdateVM(d.Id(), d.Get("denomination").(string), VM.Description, VM.CPU.Processors, VM.Memory, VM.PowerStatus); err != nil {
			log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Changing Denomination: %#v\n", err)
			return err
		}
		// d.SetPartial("denomination")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: denomination field in VM after %#v\n", VM.Denomination)
	}
	if d.HasChange("description") {
		DescriptionOldState, DescriptionNewState := d.GetChange("description")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Description %#v\n", DescriptionOldState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Description %#v\n", DescriptionNewState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Description field in VM before %#v\n", VM.Description)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, d.Get("description").(string), VM.CPU.Processors, VM.Memory, VM.PowerStatus); err != nil {
			log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Changing Description: %#v\n", err)
			return err
		}
		// d.SetPartial("denomination")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Description field in VM after %#v\n", VM.Description)
	}
	if d.HasChange("processors") {
		ProcessorsOldState, ProcessorsNewState := d.GetChange("processors")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Processors %#v\n", ProcessorsOldState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Processors %#v\n", ProcessorsNewState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Processors field in VM before %#v\n", VM.CPU.Processors)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, VM.Description, d.Get("processors").(int), VM.Memory, VM.PowerStatus); err != nil {
			log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Changing CPU %#v\n", err)
			return err
		}
		// d.SetPartial("processors")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Processors field in VM after %#v\n", VM.CPU.Processors)
	}
	if d.HasChange("memory") {
		MemoryOldState, MemoryNewState := d.GetChange("memory")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of Memory %#v\n", MemoryOldState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of Memory %#v\n", MemoryNewState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Memory field in VM before %#v\n", VM.Memory)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, VM.Description, VM.CPU.Processors, d.Get("memory").(int), VM.PowerStatus); err != nil {
			log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Changing Memory: %#v\n", err)
			return err
		}
		// d.SetPartial("memory")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Memory field in VM after %#v\n", VM.Memory)
	}
	if d.HasChange("state") {
		OldState, NewState := d.GetChange("state")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Old value of State %#v\n", OldState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: New value of State %#v\n", NewState)
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: State field in VM before %#v\n", VM.PowerStatus)
		if VM, err = apiClient.UpdateVM(d.Id(), VM.Denomination, VM.Description, VM.CPU.Processors, VM.Memory, d.Get("state").(string)); err != nil {
			log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Error Changing Memory: %#v\n", err)
			return err
		}
		// d.SetPartial("memory")
		log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj: Memory field in VM after %#v\n", VM.Memory)
	}
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmUpdate Obj:DataScheme %#v\n", d)
	d.Partial(false) // this is such as to a semaphore, it's a switch to change a state of unblocked
	return resourceVMWSVmRead(d, m)
}

func resourceVMWSVmDelete(d *schema.ResourceData, m interface{}) error {
	// We have to check if the instance is powered on
	// because before to remove we need tp make sure that
	// is shutdown
	apiClient := m.(*wsapiclient.Client)
	err := apiClient.DeleteVM(d.Id())
	if err != nil {
		log.Printf("[ERROR][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmDelete Obj:%#v\n", err)
		return err
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
	log.Printf("[DEBUG][VMWS] Fi: resource_vmworkstation_vm.go Fu: resourceVMWSVmExists Obj:APIClient %#v\n", apiClient)
	if VM == nil {
		return false, nil
	}
	return true, nil
}

func resourceVMWSVmImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
