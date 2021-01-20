package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADComputer() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_computer` manages computer objects in an AD tree.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Read:   resourceADComputerRead,
		Create: resourceADComputerCreate,
		Update: resourceADComputerUpdate,
		Delete: resourceADComputerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressCaseDiff,
				Description:      "The name for the computer account.",
			},
			"pre2kname": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressCaseDiff,
				Description:      "The pre-win2k name for the computer account.",
			},
			"container": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCaseDiff,
				Description:      "The DN of the container used to hold the computer account.",
				Computed:         true,
			},
			"dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the computer object.",
			},
		},
	}
}

func resourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	computer, err := winrmhelper.NewComputerFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ObjectNotFound") {
			// Resource no longer exists
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while reading computer with GUID %q: %s", d.Id(), err)
	}
	_ = d.Set("name", computer.Name)
	_ = d.Set("dn", computer.DN)
	_ = d.Set("guid", computer.GUID)
	_ = d.Set("pre2kname", computer.SAMAccountName)
	_ = d.Set("container", computer.Path)
	_ = d.Set("sid", computer.SID.Value)

	return nil
}

func resourceADComputerCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	computer := winrmhelper.NewComputerFromResource(d)

	guid, err := computer.Create(client)
	if err != nil {
		return fmt.Errorf("error while creating new computer object: %s", err)
	}
	d.SetId(guid)
	return resourceADComputerRead(d, meta)
}

func resourceADComputerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	computer := winrmhelper.NewComputerFromResource(d)
	keys := []string{"container"}
	changes := make(map[string]interface{})
	for _, key := range keys {
		if d.HasChange(key) {
			changes[key] = d.Get(key)
		}
	}

	err = computer.Update(client, changes)
	if err != nil {
		return fmt.Errorf("error while updating computer with id %q: %s", d.Id(), err)
	}
	return resourceADComputerRead(d, meta)
}

func resourceADComputerDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	computer := winrmhelper.NewComputerFromResource(d)
	err = computer.Delete(client)
	if err != nil {
		return fmt.Errorf("error while deleting a computer object with id %q: %s", d.Id(), err)
	}

	return resourceADComputerRead(d, meta)
}
