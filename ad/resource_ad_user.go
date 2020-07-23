package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADUser() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_user` manages User objects in an Active Directory tree.",
		Create:      resourceADUserCreate,
		Read:        resourceADUserRead,
		Update:      resourceADUserUpdate,
		Delete:      resourceADUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Display Name of an Active Directory user.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Principal Name of an Active Directory user.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pre-win2k user logon name.",
			},
			"initial_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's initial password. This will be set on creation but will *not* be enforced in subsequent plans.",
			},
			"container": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "A DN of the container object that will be holding the user.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"cannot_change_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true the user will not be allowed to change their password.",
			},
			"password_never_expires": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true the password for this user will not expire.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to false the user will be disabled.",
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
		return fmt.Errorf("while retrieving user data from host: %s", err)
	}
	err = u.DeleteUser(client)
	if err != nil {
		return fmt.Errorf("while deleting user: %s", err)
	}
	return resourceADUserRead(d, meta)
}
