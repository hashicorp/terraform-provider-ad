package ad

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceADUserCreate,
		Read:   resourceADUserRead,
		Update: resourceADUserUpdate,
		Delete: resourceADUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sam_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initial_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cannot_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password_never_expires": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceADUserCreate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetUserFromResource(d)
	client := meta.(ProviderConf).WinRMClient
	guid, err := u.NewUser(client)
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADUserRead(d, meta)
}

func resourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading ad_user resource for user with guid: %q", d.Id())
	client := meta.(ProviderConf).WinRMClient
	u, err := winrmhelper.GetUserFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			d.SetId("")
			return nil
		}
		return err
	}
	if u == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("sam_account_name", u.SAMAccountName)
	_ = d.Set("display_name", u.DisplayName)
	_ = d.Set("principal_name", u.PrincipalName)
	_ = d.Set("container", u.Container)
	_ = d.Set("enabled", u.Enabled)
	_ = d.Set("password_never_expires", u.PasswordNeverExpires)
	_ = d.Set("cannot_change_password", u.CannotChangePassword)

	return nil
}

func resourceADUserUpdate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetUserFromResource(d)
	client := meta.(ProviderConf).WinRMClient
	err := u.ModifyUser(d, client)
	if err != nil {
		return err
	}
	return resourceADUserRead(d, meta)
}

func resourceADUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	u, err := winrmhelper.GetUserFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	u.DeleteUser(client)
	return resourceADUserRead(d, meta)
}
