package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_group` manages Group objects in an Active Directory tree.",
		Create:      resourceADGroupCreate,
		Read:        resourceADGroupRead,
		Update:      resourceADGroupUpdate,
		Delete:      resourceADGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.StringInSlice([]string{"global", "domainlocal", "universal"}, false),
				Description:  "The group's scope. Can be one of `global`, `domainlocal`, or `universal` (case sensitive).",
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "security",
				ValidateFunc: validation.StringInSlice([]string{"system", "security"}, false),
				Description:  "The group's category. Can be one of `system` or `security` (case sensitive).",
			},
			"container": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "A DN of a container object holding the group.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the group object.",
			},
		},
	}
}

func resourceADGroupCreate(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	u := winrmhelper.GetGroupFromResource(d)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	guid, err := u.AddGroup(client, isLocal)
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADGroupRead(d, meta)
}

func resourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	log.Printf("Reading ad_Group resource for group with GUID: %q", d.Id())
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	g, err := winrmhelper.GetGroupFromHost(client, d.Id(), isLocal)
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
	_ = d.Set("container", g.Container)
	_ = d.Set("sid", g.SID.Value)

	return nil
}

func resourceADGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	g := winrmhelper.GetGroupFromResource(d)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	err = g.ModifyGroup(d, client, isLocal)
	if err != nil {
		return err
	}
	return resourceADGroupRead(d, meta)
}

func resourceADGroupDelete(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	conn, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(conn)

	g, err := winrmhelper.GetGroupFromHost(conn, d.Id(), isLocal)
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	err = g.DeleteGroup(conn, isLocal)
	if err != nil {
		return fmt.Errorf("while deleting group: %s", err)
	}
	return nil
}
