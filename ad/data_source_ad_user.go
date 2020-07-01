package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"

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
			"domain_dn": {
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
	conn := meta.(ProviderConf).LDAPConn
	dn := d.Get("user_dn").(string)

	u, err := ldaphelper.GetUserFromLDAP(conn, dn)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("No user found with dn %q", dn)
	}
	d.Set("sam_account_name", u.SAMAccountName)
	d.Set("display_name", u.DisplayName)
	d.Set("principal_name", u.PrincipalName)
	d.SetId(dn)
	return nil
}
