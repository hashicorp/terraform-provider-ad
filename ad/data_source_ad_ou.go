package ad

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADOU() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Organizational Unit Active Directory object.",
		Read:        dataSourceADOURead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the OU object. If this is used then the `path` attribute needs to be set as well.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the OU object. If this is used then the `Name` attribute needs to be set as well.",
			},
			"dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Distinguished Name of the OU object.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OU's description.",
			},
			"protected": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OU's protected status.",
			},
		},
	}
}

func dataSourceADOURead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

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

	d.SetId(ou.GUID)
	return nil

}
