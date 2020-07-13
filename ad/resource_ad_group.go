package ad

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"
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
			"domain_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "security",
				ValidateFunc: validation.StringInSlice([]string{"system", "security"}, false),
			},
		},
	}
}

func resourceADGroupCreate(d *schema.ResourceData, meta interface{}) error {
	u := ldaphelper.GetGroupFromResource(d)
	conn := meta.(ProviderConf).LDAPConn
	dn, err := u.AddGroup(conn)
	if err != nil {
		return err
	}
	d.SetId(*dn)
	return resourceADGroupRead(d, meta)
}

func resourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading ad_Group resource for DN: %q", d.Id())
	conn := meta.(ProviderConf).LDAPConn
	g, err := ldaphelper.GetGroupFromLDAP(conn, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "No entries found for filter") {
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
	_ = d.Set("display_name", g.Name)
	_ = d.Set("domain_dn", g.DomainDN)
	_ = d.Set("scope", g.Scope)
	_ = d.Set("type", g.Type)

	return nil
}

func resourceADGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	u := ldaphelper.GetGroupFromResource(d)
	conn := meta.(ProviderConf).LDAPConn
	err := u.ModifyGroup(d, conn)
	if err != nil {
		return err
	}
	return resourceADGroupRead(d, meta)
}

func resourceADGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(ProviderConf).LDAPConn
	delReq := ldap.NewDelRequest(d.Id(), []ldap.Control{})
	delReq.DN = d.Id()
	err := conn.Del(delReq)
	if err != nil {
		return err
	}
	return resourceADGroupRead(d, meta)
}
