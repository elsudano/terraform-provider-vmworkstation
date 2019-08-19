package vmworkstation

import (
	"fmt"
	"log"

	wsapiclient "github.com/elsudano/vmware-workstation-api-client"
	"github.com/hashicorp/terraform/helper/schema"
)

type Config struct {
	User         string
	Password     string
	URL          string
	InsecureFlag bool
	Debug        bool
}

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

	if d.Get("debug") == true {
		log.Printf("[VMWS] Fu: NewConfig Fi: config.go Ob: %#v\n", config)
	}
	return config, nil
}

func (c *Config) Client() (*wsapiclient.Client, error) {
	client, err := wsapiclient.NewClient(c.User, c.Password, c.URL)

	if err != nil {
		return nil, err
	}

	if c.Debug == true {
		log.Printf("[VMWS] Fu: Client Fi: config.go Ob: %#v\n", client)
	}
	return client, err
}
