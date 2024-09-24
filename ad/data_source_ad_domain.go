package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory Computer object.",
		Read:        dataSourceADDomainRead,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OU's identifier. It can be the OU's GUID, SID, Distinguished Name, or SAM Account Name.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The GUID of the domain object.",
			},
			"dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Distinguished Name of the domain object.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the domain object.",
				Computed:    true,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the domain object.",
			},
		},
	}
}

func dataSourceADDomainRead(d *schema.ResourceData, meta interface{}) error {
	computerID := winrmhelper.SanitiseTFInput(d, "domain_id")

	var identity string
	if computerID == "" {
		return fmt.Errorf("invalid inputs for AD computer datasource. domain_id is required")
	} else {
		identity = computerID
	}

	domain, err := winrmhelper.NewDomainFromHost(meta.(*config.ProviderConf), identity)
	if err != nil {
		return err
	}

	d.SetId(domain.GUID)
	_ = d.Set("name", domain.Name)
	_ = d.Set("dn", domain.DN)
	_ = d.Set("guid", domain.GUID)
	_ = d.Set("sid", domain.SID.Value)

	return nil
}
