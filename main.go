package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-scaffolding/msad" // change this to the import path of your provider
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: msad.Provider})
}
