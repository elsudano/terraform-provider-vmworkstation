package vmworkstation

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVMWSFolder() *schema.Resource {
	return &schema.Resource{
		Read: datasourceVMWSFolderRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Description: "The absolute path of the folder.",
				Required:    true,
			},
		},
	}
}

func datasourceVMWSFolderRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
