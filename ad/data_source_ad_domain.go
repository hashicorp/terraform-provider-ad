package ad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"
)

func dataSourceADDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADDomainRead,

		Schema: map[string]*schema.Schema{
			"netbios_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceADDomainRead(d *schema.ResourceData, meta interface{}) error {
	nb := d.Get("netbios_name").(string)
	dns := d.Get("domain_name").(string)
	dn := d.Get("domain_dn").(string)

	dseConn := meta.(ProviderConf).LDAPDSEConn
	domain, err := ldaphelper.GetDomainFromLDAP(dseConn, dn, nb, dns)
	if err != nil {
		return err
	}

	_ = d.Set("netbios_name", domain.NetbiosName)
	_ = d.Set("domain_name", domain.DomainName)
	_ = d.Set("domain_dn", domain.DN)
	d.SetId(domain.DN)

	return nil
}
