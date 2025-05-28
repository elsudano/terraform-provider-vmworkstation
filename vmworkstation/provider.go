package vmworkstation

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider method is the entry point to use the provider
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_USER", nil),
				Description: "The user name for VMWare Workstation Pro API REST operations.",
				Sensitive:   true,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_PASSWORD", nil),
				Description: "The user password for VMWare Workstation Pro API REST operations.",
				Sensitive:   true,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_URL", nil),
				Description: "The URL for connect to the API REST",
			},
			"https": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_HTTPS", true),
				Description: "When this have set to true the 'url' connect to over https",
			},
			"debug": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VMWS_DEBUG", nil),
				Description: "Enable debug for find errors",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vmworkstation_vm":     resourceVMWSVm(),
			"vmworkstation_folder": resourceVMWSFolder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"vmworkstation_folder": datasourceVMWSFolder(),
		},
		ConfigureFunc: providerConfigure,
	}
	log.Printf("[DEBUG][VMWS] Fi: provider.go Fu: Provider Obj:%#v\n", provider)
	return provider
}

// providerConfigure this method give a new configuration object to use with the client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config, err := NewConfig(d)
	if err != nil {
		return nil, err
	}
	myClient, err := config.Client()
	// myClient.SwitchDebug()
	log.Printf("[DEBUG][VMWS] Fi: provider.go Fu: providerConfigure Obj:%#v\n", d.State().String())
	log.Printf("[DEBUG][VMWS] Fi: provider.go Fu: providerConfigure Obj:%#v\n", config)
	return myClient, err
}
