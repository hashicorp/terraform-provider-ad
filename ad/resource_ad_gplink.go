package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGPLink() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_gplink` manages links between GPOs and container objects such as OUs.",
		Create:      resourceADGPLinkCreate,
		Read:        resourceADGPLinkRead,
		Update:      resourceADGPLinkUpdate,
		Delete:      resourceADGPLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gpo_guid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					_, err := uuid.ParseUUID(val.(string))
					if err != nil {
						errs = append(errs, fmt.Errorf("%q is not a valid uuid", val.(string)))
					}
					return
				},
				Description:      "The GUID of the GPO that will be linked to the container object.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"target_dn": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The DN of the object the GPO will be linked to.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"enforced": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true the GPO will be enforced on the container object.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Controls the state of the GP link between a GPO and a container object.",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the precedence between multiple GPOs linked to the same container object.",
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

	_ = d.Set("gpo_guid", gplink.GPOGuid)
	_ = d.Set("target_dn", gplink.Target)
	_ = d.Set("enforced", gplink.Enforced)
	_ = d.Set("enabled", gplink.Enabled)
	_ = d.Set("order", gplink.Order)

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
