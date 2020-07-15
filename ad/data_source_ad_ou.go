package ad

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADOU() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADOURead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceADOURead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient

	name := winrmhelper.SanitiseTFInput(d, "name")
	path := winrmhelper.SanitiseTFInput(d, "path")
	dn := winrmhelper.SanitiseTFInput(d, "dn")

	if dn == "" && (name == "" || path == "") {
		return fmt.Errorf("invalid inputs, dn or a combination of path and name are required")
	}

	ou, err := winrmhelper.NewOrgUnitFromHost(client, dn, name, path)
	if err != nil {
		return err
	}

	_ = d.Set("name", ou.Name)
	_ = d.Set("description", ou.Description)
	_ = d.Set("path", ou.Path)
	_ = d.Set("protected", strconv.FormatBool(ou.Protected))
	_ = d.Set("dn", ou.DistinguishedName)

	d.SetId(ou.DistinguishedName)
	return nil

}
