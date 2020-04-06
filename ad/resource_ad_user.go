package ad

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"
)

func resourceADUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceADUserCreate,
		Read:   resourceADUserRead,
		Update: resourceADUserUpdate,
		Delete: resourceADUserDelete,

		Schema: map[string]*schema.Schema{
			"domain_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
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
		},
	}
}

func resourceADUserCreate(d *schema.ResourceData, meta interface{}) error {
	u := ldaphelper.GetUserFromResource(d)
	conn := meta.(ProviderConf).LDAPConn
	dn, err := u.AddUser(conn)
	if err != nil {
		return err
	}
	d.SetId(*dn)
	return resourceADUserRead(d, meta)
}

func resourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading ad_user resource for DN: %q", d.Id())
	conn := meta.(ProviderConf).LDAPConn
	domainDN := d.Get("domain_dn").(string)
	u, err := ldaphelper.GetUserFromLDAP(conn, d.Id(), domainDN)
	if err != nil {
		if strings.Contains(err.Error(), "No entries found for filter") {
			d.SetId("")
			return nil
		}
		return err
	}
	if u == nil {
		d.SetId("")
		return nil
	}
	d.Set("sam_account_name", u.SAMAccountName)
	d.Set("display_name", u.DisplayName)
	d.Set("principal_name", u.PrincipalName)

	return nil
}

func resourceADUserUpdate(d *schema.ResourceData, meta interface{}) error {
	u := ldaphelper.GetUserFromResource(d)
	conn := meta.(ProviderConf).LDAPConn
	err := u.ModifyUser(d, conn)
	if err != nil {
		return err
	}
	return resourceADUserRead(d, meta)
}

func resourceADUserDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(ProviderConf).LDAPConn
	delReq := ldap.NewDelRequest(d.Id(), []ldap.Control{})
	delReq.DN = d.Id()
	err := conn.Del(delReq)
	if err != nil {
		return err
	}
	return resourceADUserRead(d, meta)
}
