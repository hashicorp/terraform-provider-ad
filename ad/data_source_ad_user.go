package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceADUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADUserRead,

		Schema: map[string]*schema.Schema{
			"user_dn": {
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
			"principal_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	dn := d.Get("user_dn").(string)
	client := meta.(ProviderConf).WinRMClient
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
	d.SetId(u.GUID)

	return nil
}
