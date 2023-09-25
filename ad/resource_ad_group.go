// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

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
				ValidateFunc: validation.StringInSlice([]string{"distribution", "security"}, false),
				Description:  "The group's category. Can be one of `distribution` or `security` (case sensitive).",
			},
			"container": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "A DN of a container object holding the group.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Group.",
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the group object.",
			},
			"dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name of the group object.",
			},
		},
	}
}

func resourceADGroupCreate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetGroupFromResource(d)
	guid, err := u.AddGroup(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADGroupRead(d, meta)
}

func resourceADGroupRead(d *schema.ResourceData, meta interface{}) error {
	g, err := winrmhelper.GetGroupFromHost(meta.(*config.ProviderConf), d.Id())
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
	_ = d.Set("description", g.Description)
	_ = d.Set("dn", g.DistinguishedName)
	_ = d.Set("sid", g.SID.Value)

	return nil
}

func resourceADGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	g := winrmhelper.GetGroupFromResource(d)
	err := g.ModifyGroup(d, meta.(*config.ProviderConf))
	if err != nil {
		return err
	}
	return resourceADGroupRead(d, meta)
}

func resourceADGroupDelete(d *schema.ResourceData, meta interface{}) error {
	g, err := winrmhelper.GetGroupFromHost(meta.(*config.ProviderConf), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	err = g.DeleteGroup(meta.(*config.ProviderConf))
	if err != nil {
		return fmt.Errorf("while deleting group: %s", err)
	}
	return nil
}
