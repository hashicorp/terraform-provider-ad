package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceADGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory Group object.",
		Read:        dataSourceADGroupRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the Group object. Alternatively it can be the SID, the Distinguished Name, or the SAM Account Name of the group.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAM account name of the Group object.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the Group object.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Group object.",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Group's category.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Group's scope.",
			},
			"container": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Group's container object.",
			},
		},
	}
}

func dataSourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	dn := d.Get("guid").(string)

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
	_ = d.Set("name", g.Name)

	d.SetId(g.GUID)
	return nil
}
