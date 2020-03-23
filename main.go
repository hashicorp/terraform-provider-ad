package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-msad/msad"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: msad.Provider})
}
