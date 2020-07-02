package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceADGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADGroupRead,

		Schema: map[string]*schema.Schema{
			"group_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sam_account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	dn := d.Get("group_dn").(string)

	g, err := winrmhelper.GetGroupFromHost(client, dn)
	if err != nil {
		return err
	}
	if g == nil {
		return fmt.Errorf("No group found with dn %q", dn)
	}
	_ = d.Set("sam_account_name", g.SAMAccountName)
	_ = d.Set("display_name", g.Name)
	_ = d.Set("scope", g.Scope)
	_ = d.Set("category", g.Category)
	_ = d.Set("container", g.Container)

	d.SetId(dn)
	return nil
}
