package vmworkstation

import (
	"fmt"
	"log"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform/helper/schema"
)

// Config this struct is to storage the data of configuration to the web client
type Config struct {
	User         string
	Password     string
	URL          string
	InsecureFlag bool
	Debug        bool
}

// NewConfig_bak this method is usefull tu create a new config for the client without parameters
func NewConfig(d *schema.ResourceData) (*Config, error) {
	if d.Get("user").(string) == "" || d.Get("password").(string) == "" || d.Get("url").(string) == "" {
		err := fmt.Errorf("User, Password and URL that required parameters")
		return nil, err
	}
	config := &Config{
		User:         d.Get("user").(string),
		Password:     d.Get("password").(string),
		URL:          d.Get("url").(string),
		InsecureFlag: d.Get("https").(bool),
		Debug:        d.Get("debug").(bool),
	}
	// log.Printf("[VMWS] Fi: config.go Fu: NewConfig Obj:%#v\n", config)
	return config, nil
}

// Client_bak this method give a new client to communicate with the API REST to VmWare Workstation
func (co *Config) Client() (*wsapiclient.Client, error) {
	client, err := wsapiclient.New()
	// log.Printf("[VMWS] Fi: config.go Fu: Client Obj: client after %#v\n", client)
	client.ConfigCli(co.URL, co.User, co.Password, co.InsecureFlag, co.Debug)
	if err != nil {
		return nil, err
	}
	log.Printf("[VMWS] Fi: config.go Fu: Client Obj: client before %#v\n", client)
	return client, err
}
