package ad

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceADGroupCreate,
		Read:   resourceADGroupRead,
		Update: resourceADGroupUpdate,
		Delete: resourceADGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sam_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "global",
				ValidateFunc: validation.StringInSlice([]string{"global", "local", "universal"}, false),
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "security",
				ValidateFunc: validation.StringInSlice([]string{"system", "security"}, false),
			},
			"container": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceADGroupCreate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetGroupFromResource(d)
	client := meta.(ProviderConf).WinRMClient
	guid, err := u.AddGroup(client)
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADGroupRead(d, meta)
}

func resourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading ad_Group resource for group with GUID: %q", d.Id())
	client := meta.(ProviderConf).WinRMClient
	g, err := winrmhelper.GetGroupFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			d.SetId("")
			return nil
		}
		return err
	}
	if g == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("sam_account_name", g.SAMAccountName)
	_ = d.Set("name", g.Name)
	_ = d.Set("scope", g.Scope)
	_ = d.Set("category", g.Category)

	return nil
}

func resourceADGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	g := winrmhelper.GetGroupFromResource(d)
	client := meta.(ProviderConf).WinRMClient
	err := g.ModifyGroup(d, client)
	if err != nil {
		return err
	}
	return resourceADGroupRead(d, meta)
}

func resourceADGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(ProviderConf).WinRMClient
	g, err := winrmhelper.GetGroupFromHost(conn, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	g.DeleteGroup(conn)
	return resourceADGroupRead(d, meta)
}
