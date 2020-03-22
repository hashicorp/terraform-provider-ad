package msad

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type ProviderMeta struct {
	LDAPClient interface{}
}

// Provider exports the provider schema
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bind_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_USER", nil),
				Description: "The username used to authenticate to the AD's LDAP service.",
			},
			"bind_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PASSWORD", nil),
				Description: "The username used to authenticate to the AD's LDAP service.",
			},
			"dc_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_HOSTNAME", nil),
				Description: "The username used to authenticate to the AD's LDAP service.",
			},
			"dc_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PORT", 389),
				Description: "The username used to authenticate to the AD's LDAP service.",
			},
			"allow_insecure_certs": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_INSECURE", false),
				Description: "The username used to authenticate to the AD's LDAP service.",
			},
			"proto": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PROTO", "ldap"),
				Description: "The protocol to use when talking to AD. Valid choices are ldap or ldaps",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"msad_domain": dataSourceMSADDomain(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// "scaffolding_resource": resourceScaffolding(),
		},
		ConfigureFunc: getProviderConfig,
	}
}

func getProviderConfig(d *schema.ResourceData) (interface{}, error) {

	return nil, nil
}
