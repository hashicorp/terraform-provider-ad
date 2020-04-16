package ad

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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
				Type:        schema.TypeBool,
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
			"ad_domain": dataSourceADDomain(),
			"ad_user":   dataSourceADUser(),
			"ad_group":  dataSourceADGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ad_user":  resourceADUser(),
			"ad_group": resourceADGroup(),
		},
		ConfigureFunc: initProviderConfig,
	}
}

// ProviderConf holds structures that are useful to the provider at runtime.
type ProviderConf struct {
	Configuration *ProviderConfig
	LDAPConn      *ldap.Conn
	LDAPDSEConn   *ldap.Conn
}

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {

	cfg := NewConfig(d)
	conn, err := GetConnection(cfg, false)
	if err != nil {
		return nil, err
	}

	rootDseConn, err := GetConnection(cfg, true)
	if err != nil {
		return nil, err
	}

	pcfg := ProviderConf{
		Configuration: &cfg,
		LDAPConn:      conn,
		LDAPDSEConn:   rootDseConn,
	}

	return pcfg, nil
}
