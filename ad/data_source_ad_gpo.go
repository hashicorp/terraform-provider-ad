package ad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADGPO() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADGPORead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceADGPORead(d *schema.ResourceData, meta interface{}) error {
	name := winrmhelper.SanitiseTFInput(d, "name")
	guid := winrmhelper.SanitiseTFInput(d, "guid")

	client := meta.(ProviderConf).WinRMClient

	gpo, err := winrmhelper.GetGPOFromHost(client, name, guid)
	if err != nil {
		return err
	}

	d.Set("name", gpo.Name)
	d.Set("domain", gpo.Domain)
	d.SetId(gpo.ID)

	return nil
}
