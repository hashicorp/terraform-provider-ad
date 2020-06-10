package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceADGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADGroupRead,

		Schema: map[string]*schema.Schema{
			"dn": {
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(ProviderConf).LDAPConn
	dn := d.Get("dn").(string)
	domainDN := d.Get("domain_dn").(string)

	g, err := ldaphelper.GetGroupFromLDAP(conn, dn, domainDN)
	if err != nil {
		return err
	}
	if g == nil {
		return fmt.Errorf("No group found with dn %q", dn)
	}
	d.Set("sam_account_name", g.SAMAccountName)
	d.Set("display_name", g.Name)
	d.Set("scope", g.Scope)
	d.Set("type", g.Type)

	d.SetId(dn)
	return nil
}
