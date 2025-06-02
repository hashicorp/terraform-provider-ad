// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ad

import (
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADOU() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_ou` manages OU objects in an AD tree.",
		Read:        resourceADOURead,
		Create:      resourceADOUCreate,
		Update:      resourceADOUUpdate,
		Delete:      resourceADOUDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the OU.",
			},
			"path": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "DN of the object that contains the OU.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the OU.",
			},
			"protected": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Protect this OU from being deleted accidentaly.",
			},
			"dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OU's DN.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OU's GUID.",
			},
		},
	}
}

func resourceADOURead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	ou, err := winrmhelper.NewOrgUnitFromHost(meta.(*config.ProviderConf), d.Id(), "", "")
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
	_ = d.Set("path", ou.Path)
	_ = d.Set("protected", ou.Protected)
	_ = d.Set("dn", ou.DistinguishedName)
	_ = d.Set("guid", ou.GUID)

	return nil
}

func resourceADOUCreate(d *schema.ResourceData, meta interface{}) error {
	ou := winrmhelper.NewOrgUnitFromResource(d)
	guid, err := ou.Create(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}
	d.SetId(guid)

	return resourceADOURead(d, meta)
}

func resourceADOUUpdate(d *schema.ResourceData, meta interface{}) error {
	ou := winrmhelper.NewOrgUnitFromResource(d)

	keys := []string{"description", "name", "path", "protected"}
	changes := make(map[string]interface{})
	for _, key := range keys {
		if d.HasChange(key) {
			changes[key] = d.Get(key)
		}
	}

	err := ou.Update(meta.(*config.ProviderConf), changes)
	if err != nil {
		return err
	}
	return resourceADOURead(d, meta)
}

func resourceADOUDelete(d *schema.ResourceData, meta interface{}) error {
	ou := winrmhelper.NewOrgUnitFromResource(d)
	err := ou.Delete(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}

	return nil
}
