package vmworkstation

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_USER", nil),
				Description: "The user name for VMWare Workstation Pro API REST operations.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_PASSWORD", nil),
				Description: "The user password for VMWare Workstation Pro API REST operations.",
			},

			"url_to_api": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_URL_API", nil),
				Description: "The URL for connect to the API REST",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmws_vms": resourceVMWSVms(),
			"vmws_folder": resourceVMWSFolder(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vmws_folder": datasourceVMWSFolder(),
		},

		// ConfigureFunc: providerConfigure,
	}
}

/* func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := NewConfig(d)
	if err != nil {
		return nil, err
	}
	return c.Client()
} */
