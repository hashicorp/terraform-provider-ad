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
		Description: "`ad_group` manages Group objects in an Active directory tree.",
		Create:      resourceADGroupCreate,
		Read:        resourceADGroupRead,
		Update:      resourceADGroupUpdate,
		Delete:      resourceADGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pre-win2k name of the group.",
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "global",
				ValidateFunc: validation.StringInSlice([]string{"global", "local", "universal"}, false),
				Description:  "The group's scope. Can be one of `global`, `local`, or `universal` (case sensitive).",
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "security",
				ValidateFunc: validation.StringInSlice([]string{"system", "security"}, false),
				Description:  "The group's category. Can be one of `system` or `security` (case sensitive).",
			},
			"container": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A DN of a container object holding the group.",
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
