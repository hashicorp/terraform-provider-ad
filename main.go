package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-ad/ad"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ad.Provider})
}
