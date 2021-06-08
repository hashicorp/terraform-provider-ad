package ad

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADOU() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Organizational Unit Active Directory object.",
		Read:        dataSourceADOURead,
		Schema: map[string]*schema.Schema{
			"ou_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OU's identifier. It can be the OU's GUID, SID, Distinguished Name, or SAM Account Name.",
			},
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
				Deprecated:  "This field is deprecated in favour of `ou_id`. In the future this field will be read-only. This field is deprecated in favour of `computer_id`. In the future this field will be read-only.",
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
	name := winrmhelper.SanitiseTFInput(d, "name")
	path := winrmhelper.SanitiseTFInput(d, "path")
	dn := winrmhelper.SanitiseTFInput(d, "dn")
	ouID := winrmhelper.SanitiseTFInput(d, "ou_id")

	if dn == "" && (name == "" || path == "") && ouID == "" {
		return fmt.Errorf("invalid inputs, ou_id or dn or a combination of path and name are required")
	}

	var ouIdentifier string
	if ouID != "" {
		ouIdentifier = ouID
	} else {
		ouIdentifier = dn
	}
	ou, err := winrmhelper.NewOrgUnitFromHost(meta.(*config.ProviderConf), ouIdentifier, name, path)
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
