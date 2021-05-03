package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the computer object.",
			},
		},
	}
}

func dataSourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	isPassCredentialsEnabled := meta.(ProviderConf).isPassCredentialsEnabled()
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

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

	computer, err := winrmhelper.NewComputerFromHost(client, identity, isLocal, isPassCredentialsEnabled, meta.(ProviderConf).Configuration.WinRMUsername, meta.(ProviderConf).Configuration.WinRMPassword)
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
