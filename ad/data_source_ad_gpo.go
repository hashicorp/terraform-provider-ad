package ad

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADGPO() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory Group Policy Object.",
		Read:        dataSourceADGPORead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the GPO.",
				Optional:    true,
			},
			"guid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GUID of the GPO.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain of the GPO.",
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
		if strings.Contains(err.Error(), "GpoWithNameNotFound") || strings.Contains(err.Error(), "GpoWithIdNotFound") {
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("name", gpo.Name)
	_ = d.Set("domain", gpo.Domain)
	d.SetId(gpo.ID)

	return nil
}
