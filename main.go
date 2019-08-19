package main

import (
	"github.com/elsudano/terraform-provider-vmworkstation/vmworkstation"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vmworkstation.Provider,
	})
}

/* func main() {
	client, err := wsapiclient.New()
	if err != nil {
		log.Printf("[VMWS] Fu: Main Fi: main.go Ob: %#v\n", client)
		os.Exit(1)
	}
	fmt.Printf("La URL: %s\n", client.BaseURL.String())
	client.User = "admin"
	client.Password = "Carfer#007"
	client.BaseURL, _ = url.Parse("https://localhost:5555/api/")
	fmt.Printf("La URL: %s\n", client.BaseURL.String())
} */
