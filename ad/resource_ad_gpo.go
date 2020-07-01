package ad

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGPO() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_gpo` manages Group Policy Objects (GPOs).",
		Create:      resourceADGPOCreate,
		Read:        resourceADGPORead,
		Update:      resourceADGPOUpdate,
		Delete:      resourceADGPODelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Group Policy Object.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain of the GPO.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the GPO.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AllSettingsEnabled",
				ValidateFunc: validation.StringInSlice([]string{"AllSettingsEnabled", "UserSettingsDisabled", "ComputerSettingsDisabled", "AllSettingsDisabled"}, false),
				Description:  "Status of the GPO. Can be one of `AllSettingsEnabled`, `UserSettingsDisabled`, `ComputerSettingsDisabled`, or `AllSettingsDisabled` (case sensitive).",
			},
			"numeric_status": {
				Type:     schema.TypeInt,
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
	_ = d.Set("domain", g.Domain)
	_ = d.Set("description", g.Description)
	_ = d.Set("status", g.Status)
	_ = d.Set("numeric_status", g.NumericStatus)
	_ = d.Set("name", g.Name)
	return nil
}

func resourceADGPOUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	g := winrmhelper.GetGPOFromResource(d)
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
