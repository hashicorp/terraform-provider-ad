package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

func resourceADComputer() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			},
			"pre2kname": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressCaseDiff,
			},
			"container": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "CN=Computers,DC=yourdomain,DC=com",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	client := meta.(ProviderConf).WinRMClient

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

	return nil
}

func resourceADComputerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	computer := winrmhelper.NewComputerFromResource(d)

	guid, err := computer.Create(client)
	if err != nil {
		return fmt.Errorf("error while creating new computer object: %s", err)
	}
	d.SetId(guid)
	return resourceADComputerRead(d, meta)
}

func resourceADComputerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	computer := winrmhelper.NewComputerFromResource(d)
	keys := []string{"container"}
	changes := make(map[string]interface{})
	for _, key := range keys {
		if d.HasChange(key) {
			changes[key] = d.Get(key)
		}
	}

	err := computer.Update(client, changes)
	if err != nil {
		return fmt.Errorf("error while updating computer with id %q: %s", d.Id(), err)
	}
	return resourceADComputerRead(d, meta)
}

func resourceADComputerDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}
	client := meta.(ProviderConf).WinRMClient
	computer := winrmhelper.NewComputerFromResource(d)
	err := computer.Delete(client)
	if err != nil {
		return fmt.Errorf("error while deleting a computer object with id %q: %s", d.Id(), err)
	}

	return resourceADComputerRead(d, meta)
}
