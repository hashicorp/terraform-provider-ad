package msad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMSADDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSADDomainRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceMSADDomainRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
