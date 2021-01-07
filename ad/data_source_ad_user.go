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
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.",
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
	userID := d.Get("user_id").(string)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, userID)
	if err != nil {
		return err
	}

	if u == nil {
		return fmt.Errorf("No user found with user_id %q", userID)
	}
	_ = d.Set("sam_account_name", u.SAMAccountName)
	_ = d.Set("display_name", u.DisplayName)
	_ = d.Set("principal_name", u.PrincipalName)
	_ = d.Set("user_id", userID)
	d.SetId(u.GUID)

	return nil
}
