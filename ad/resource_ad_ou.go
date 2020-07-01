package ad

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADOU() *schema.Resource {
	return &schema.Resource{
		Read:   resourceADOURead,
		Create: resourceADOUCreate,
		Update: resourceADOUUpdate,
		Delete: resourceADOUDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceADOURead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}
	client := meta.(ProviderConf).WinRMClient

	ou, err := winrmhelper.NewOrgUnitFromHost(client, d.Id(), "", "")
	if err != nil {
		if strings.Contains(err.Error(), "ObjectNotFound") {
			// Resource no longer exists
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("name", ou.Name)
	_ = d.Set("description", ou.Description)
	_ = d.Set("path", strings.ToLower(ou.Path))
	_ = d.Set("protected", ou.Protected)
	_ = d.Set("dn", strings.ToLower(ou.DistinguishedName))
	_ = d.Set("guid", ou.GUID)

	return nil
}

func resourceADOUCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	ou := winrmhelper.NewOrgUnitFromResource(d)
	guid, err := ou.Create(client)
	if err != nil {
		return err
	}
	d.SetId(guid)

	return resourceADOURead(d, meta)
}

func resourceADOUUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	ou := winrmhelper.NewOrgUnitFromResource(d)

	keys := []string{"description", "name", "path", "protected"}
	changes := make(map[string]interface{})
	for _, key := range keys {
		if d.HasChange(key) {
			changes[key] = d.Get(key)
		}
	}

	err := ou.Update(client, changes)
	if err != nil {
		return err
	}
	return resourceADOURead(d, meta)
}

func resourceADOUDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	client := meta.(ProviderConf).WinRMClient
	ou := winrmhelper.NewOrgUnitFromResource(d)
	err := ou.Delete(client)
	if err != nil {
		return err
	}

	return resourceADOURead(d, meta)
}
