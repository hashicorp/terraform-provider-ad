package ad

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGPO() *schema.Resource {
	return &schema.Resource{
		Create: resourceADGPOCreate,
		Read:   resourceADGPORead,
		Update: resourceADGPOUpdate,
		Delete: resourceADGPODelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"numeric_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceADGPOCreate(d *schema.ResourceData, meta interface{}) error {
	g := winrmhelper.GetGPOFromResource(d)
	client := meta.(ProviderConf).WinRMClient
	guid, err := g.NewGPO(client)
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADGPORead(d, meta)
}

func resourceADGPORead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}
	client := meta.(ProviderConf).WinRMClient
	g, err := winrmhelper.GetGPOFromHost(client, "", d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "GpoWithNameNotFound") || strings.Contains(err.Error(), "GpoWithIdNotFound") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("domain", g.Domain)
	d.Set("description", g.Description)
	d.Set("status", g.Status)
	d.Set("name", g.Name)
	return nil
}

func resourceADGPOUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	g := winrmhelper.GetGPOFromResource(d)
	log.Printf("[DEBUG] TTTTT %#v", g)
	_, err := g.UpdateGPO(client, d)
	if err != nil {
		return err
	}
	return resourceADGPORead(d, meta)
}

func resourceADGPODelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	g := winrmhelper.GetGPOFromResource(d)
	err := g.DeleteGPO(client)
	if err != nil {
		return err
	}
	return resourceADGPORead(d, meta)
}
