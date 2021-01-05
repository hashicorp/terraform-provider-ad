package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceADUser() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory user object.",
		Read:        dataSourceADUserRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the user object. Alternatively it can be the SID, the Distinguished Name, or the SAM Account Name of the user.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAM account name of the user object.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user object.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The principal name of the user object.",
			},
		},
	}
}

func dataSourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	dn := d.Get("guid").(string)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, dn)
	if err != nil {
		return err
	}

	if u == nil {
		return fmt.Errorf("No user found with dn %q", dn)
	}
	_ = d.Set("sam_account_name", u.SAMAccountName)
	_ = d.Set("display_name", u.DisplayName)
	_ = d.Set("principal_name", u.PrincipalName)
	_ = d.Set("guid", u.GUID)
	d.SetId(u.GUID)

	return nil
}
