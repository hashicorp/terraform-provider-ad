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
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The group's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAM account name of the Group object.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Group object.",
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
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the group object.",
			},
		},
	}
}

func dataSourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	groupID := d.Get("group_id").(string)

	g, err := winrmhelper.GetGroupFromHost(client, groupID, isLocal)
	if err != nil {
		return err
	}
	if g == nil {
		return fmt.Errorf("No group found with group_id %q", groupID)
	}
	_ = d.Set("sam_account_name", g.SAMAccountName)
	_ = d.Set("display_name", g.Name)
	_ = d.Set("scope", g.Scope)
	_ = d.Set("category", g.Category)
	_ = d.Set("container", g.Container)
	_ = d.Set("name", g.Name)
	_ = d.Set("group_id", groupID)
	_ = d.Set("description", g.Description)
	_ = d.Set("sid", g.SID.Value)

	d.SetId(g.GUID)
	return nil
}
