package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADComputer() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory Computer object.",
		Read:        dataSourceADComputerRead,
		Schema: map[string]*schema.Schema{
			"computer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OU's identifier. It can be the OU's GUID, SID, Distinguished Name, or SAM Account Name.",
			},
			"guid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The GUID of the computer object. This field is deprecated in favour of `computer_id`. In the future this field will be read-only.",
				Deprecated:  "This field is deprecated in favour of `computer_id`. In the future this field will be read-only.",
			},
			"dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Distinguished Name of the computer object. This field is deprecated in favour of `computer_id`. In the future this field will be read-only.",
				Deprecated:  "This field is deprecated in favour of `computer_id`. In the future this field will be read-only.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the computer object.",
				Computed:    true,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the computer object.",
			},
		},
	}
}

func dataSourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	dn := winrmhelper.SanitiseTFInput(d, "dn")
	guid := winrmhelper.SanitiseTFInput(d, "guid")
	computerID := winrmhelper.SanitiseTFInput(d, "computer_id")

	var identity string
	if computerID == "" && guid == "" && dn == "" {
		return fmt.Errorf("invalid inputs for AD computer datasource. computer_id dn or guid is required")
	} else if computerID != "" {
		identity = computerID
	} else if guid != "" {
		identity = guid
	} else if dn != "" {
		identity = dn
	}

	computer, err := winrmhelper.NewComputerFromHost(meta.(*config.ProviderConf), identity)
	if err != nil {
		return err
	}

	d.SetId(computer.GUID)
	_ = d.Set("name", computer.Name)
	_ = d.Set("dn", computer.DN)
	_ = d.Set("guid", computer.GUID)
	_ = d.Set("sid", computer.SID.Value)

	return nil
}
