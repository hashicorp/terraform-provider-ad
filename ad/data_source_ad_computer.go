package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func dataSourceADComputer() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory Computer object.",
		Read:        dataSourceADComputerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the computer object.",
				Computed:    true,
			},
			"guid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The GUID of the computer object.",
			},
			"dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Distinguished Name of the computer object.",
			},
		},
	}
}

func dataSourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	dn := winrmhelper.SanitiseTFInput(d, "dn")
	guid := winrmhelper.SanitiseTFInput(d, "guid")

	var identity string
	if guid == "" && dn == "" {
		return fmt.Errorf("invalid inputs for AD computer datasource. dn or guid is required")
	} else if guid != "" {
		identity = guid
	} else if dn != "" {
		identity = dn
	}

	computer, err := winrmhelper.NewComputerFromHost(client, identity)
	if err != nil {
		return err
	}

	d.SetId(computer.GUID)
	_ = d.Set("name", computer.Name)
	_ = d.Set("dn", computer.DN)
	_ = d.Set("guid", computer.GUID)

	return nil
}
