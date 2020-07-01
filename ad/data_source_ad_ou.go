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

	err = d.Set("name", ou.Name)
	if err != nil {
		return fmt.Errorf("error setting key %q to value %#v: %s", "name", ou.Name, err)
	}

	err = d.Set("description", ou.Description)
	if err != nil {
		return fmt.Errorf("error setting key %q to value %#v: %s", "description", ou.Description, err)
	}

	err = d.Set("path", ou.Path)
	if err != nil {
		return fmt.Errorf("error setting key %q to value %#v: %s", "path", ou.Path, err)
	}

	err = d.Set("protected", strconv.FormatBool(ou.Protected))
	if err != nil {
		return fmt.Errorf("error setting key %q to value %#v: %s", "protected", ou.Protected, err)
	}

	err = d.Set("dn", ou.DistinguishedName)
	if err != nil {
		return fmt.Errorf("error setting key %q to value %#v: %s", "dn", ou.DistinguishedName, err)
	}

	d.SetId(ou.DistinguishedName)
	return nil

}
