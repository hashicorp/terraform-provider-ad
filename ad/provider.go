package ad

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
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
			"winrm_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_USER", nil),
				Description: "The username used to authenticate to the the server's WinRM service.",
			},
			"winrm_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_PASSWORD", nil),
				Description: "The password used to authenticate to the the server's WinRM service.",
			},
			"winrm_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_HOSTNAME", nil),
				Description: "The hostname of the server we will use to run powershell scripts over WinRM.",
			},
			"winrm_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_PORT", 5985),
				Description: "The port WinRM is listening for connections. (default: 5985)",
			},
			"winrm_proto": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_PROTO", "http"),
				Description: "The WinRM protocol we will use. (default: http)",
			},
			"winrm_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_INSECURE", false),
				Description: "Trust unknown certificates. (default: false)",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ad_domain":   dataSourceADDomain(),
			"ad_user":     dataSourceADUser(),
			"ad_group":    dataSourceADGroup(),
			"ad_gpo":      dataSourceADGPO(),
			"ad_computer": dataSourceADComputer(),
			"ad_ou":       dataSourceADOU(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ad_user":         resourceADUser(),
			"ad_group":        resourceADGroup(),
			"ad_gpo":          resourceADGPO(),
			"ad_gpo_security": resourceADGPOSecurity(),
			"ad_computer":     resourceADComputer(),
			"ad_ou":           resourceADOU(),
		},
		ConfigureFunc: initProviderConfig,
	}
}

// ProviderConf holds structures that are useful to the provider at runtime.
type ProviderConf struct {
	Configuration *ProviderConfig
	LDAPConn      *ldap.Conn
	LDAPDSEConn   *ldap.Conn
	WinRMClient   *winrm.Client
	WinRMCPClient *winrmcp.Winrmcp
}

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {

	cfg := NewConfig(d)
	conn, err := GetLDAPConnection(cfg, false)
	if err != nil {
		return nil, err
	}

	rootDseConn, err := GetLDAPConnection(cfg, true)
	if err != nil {
		return nil, err
	}

	winRMClient, err := GetWinRMConnection(cfg)
	if err != nil {
		return nil, err
	}

	winRMCPClient, err := GetWinRMCPConnection(cfg)
	if err != nil {
		return nil, err
	}

	pcfg := ProviderConf{
		Configuration: &cfg,
		LDAPConn:      conn,
		LDAPDSEConn:   rootDseConn,
		WinRMClient:   winRMClient,
		WinRMCPClient: winRMCPClient,
	}

	return pcfg, nil
}
