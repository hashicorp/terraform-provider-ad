package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGPLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceADGPLinkCreate,
		Read:   resourceADGPLinkRead,
		Update: resourceADGPLinkUpdate,
		Delete: resourceADGPLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gpo_guid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enforced": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceADGPLinkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	idParts := strings.SplitN(d.Id(), "_", 2)
	if len(idParts) != 2 {
		return fmt.Errorf("malformed ID for GPLink resource with ID %q", d.Id())
	}
	gplink, err := winrmhelper.GetGPLinkFromHost(client, idParts[0], idParts[1])
	if err != nil {
		if strings.Contains(err.Error(), "did not find") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("while reading resource with id %q: %s", d.Id(), err)
	}

	d.Set("gpo_guid", gplink.GPOGuid)
	d.Set("target_dn", strings.ToLower(gplink.Target))
	d.Set("enforced", gplink.Enforced)
	d.Set("enabled", gplink.Enabled)
	d.Set("order", gplink.Order)

	return nil
}

func resourceADGPLinkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	gplink := winrmhelper.GetGPLinkFromResource(d)
	gpLinkID, err := gplink.NewGPLink(client)
	if err != nil {
		return fmt.Errorf("while creating GPLink resource: %s", err)
	}
	d.SetId(gpLinkID)

	return resourceADGPLinkRead(d, meta)
}

func resourceADGPLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	keys := []string{"enforced", "enabled", "order"}
	changes := make(map[string]interface{})
	for _, key := range keys {
		if d.HasChange(key) {
			changes[key] = d.Get(key)
		}
	}
	gplink := winrmhelper.GetGPLinkFromResource(d)
	err := gplink.ModifyGPLink(client, changes)
	if err != nil {
		return fmt.Errorf("while modifying GPLink with id %q: %s", d.Id(), err)
	}

	return resourceADGPLinkRead(d, meta)
}

func resourceADGPLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	gplink := winrmhelper.GetGPLinkFromResource(d)
	err := gplink.RemoveGPLink(client)
	if err != nil {
		return fmt.Errorf("while deleting resource with ID %q: %s", d.Id(), err)
	}

	return nil
}
