package msad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSADUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSADUserCreate,
		Read:   resourceMSADUserRead,
		Update: resourceMSADUserUpdate,
		Delete: resourceMSADUserDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMSADUserCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMSADUserRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMSADUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMSADUserDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
