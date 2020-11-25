package ad

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
)

// Provider exports the provider schema
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"winrm_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_USER", nil),
				Description: "The username used to authenticate to the server's WinRM service. (Environment variable: AD_USER)",
			},
			"winrm_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PASSWORD", nil),
				Description: "The password used to authenticate to the server's WinRM service. (Environment variable: AD_PASSWORD)",
			},
			"winrm_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_HOSTNAME", nil),
				Description: "The hostname of the server we will use to run powershell scripts over WinRM. (Environment variable: AD_HOSTNAME)",
			},
			"winrm_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PORT", 5985),
				Description: "The port WinRM is listening for connections. (default: 5985, environment variable: AD_PORT)",
			},
			"winrm_proto": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_PROTO", "http"),
				Description: "The WinRM protocol we will use. (default: http, environment variable: AD_PROTO)",
			},
			"winrm_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_INSECURE", false),
				Description: "Trust unknown certificates. (default: false, environment variable: AD_WINRM_INSECURE)",
			},
			"krb_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_KRB_REALM", ""),
				Description: "The name of the kerberos realm (domain) we will use for authentication. (default: \"\", environment variable: AD_KRB_REALM)",
			},
			"krb_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_KRB_CONF", ""),
				Description: "Path to kerberos configuration file. (default: none, environment variable: AD_KRB_CONF)",
			},
			"krb_spn": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_KRB_SPN", ""),
				Description: "Alternative Service Principal Name. (default: none, environment variable: AD_KRB_SPN)",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ad_user":     dataSourceADUser(),
			"ad_group":    dataSourceADGroup(),
			"ad_gpo":      dataSourceADGPO(),
			"ad_computer": dataSourceADComputer(),
			"ad_ou":       dataSourceADOU(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ad_user":             resourceADUser(),
			"ad_group":            resourceADGroup(),
			"ad_group_membership": resourceADGroupMembership(),
			"ad_gpo":              resourceADGPO(),
			"ad_gpo_security":     resourceADGPOSecurity(),
			"ad_computer":         resourceADComputer(),
			"ad_ou":               resourceADOU(),
			"ad_gplink":           resourceADGPLink(),
		},
		ConfigureFunc: initProviderConfig,
	}
}

// ProviderConf holds structures that are useful to the provider at runtime.
type ProviderConf struct {
	Configuration *ProviderConfig
	WinRMClient   *winrm.Client
	WinRMCPClient *winrmcp.Winrmcp
}

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {

	cfg := NewConfig(d)

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
		WinRMClient:   winRMClient,
		WinRMCPClient: winRMCPClient,
	}

	return pcfg, nil
}

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	// k is ignored here, but wee need to include it in the function's
	// signature in order to match the one defined for DiffSuppressFunc
	return strings.EqualFold(old, new)
}
