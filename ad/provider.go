package ad

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				//lintignore: V013
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					os := runtime.GOOS
					if v == "" && os != "windows" {
						errs = append(errs, fmt.Errorf("%q is allowed to be empty only if terraform runs on windows, (curent os: %q) ", key, os))
					}
					return
				},
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
				//lintignore: V013
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					os := runtime.GOOS
					if v == "" && os != "windows" {
						errs = append(errs, fmt.Errorf("%q is allowed to be empty only if terraform runs on windows, (curent os: %q) ", key, os))
					}
					return
				},
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
			"krb_keytab": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_KRB_KEYTAB", ""),
				Description: "Path to a keytab file to be used instead of a password",
			},
			"winrm_use_ntlm": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_USE_NTLM", false),
				Description: "Use NTLM authentication. (default: false, environment variable: AD_WINRM_USE_NTLM)",
			},
			"winrm_pass_credentials": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_PASS_CREDENTIALS", false),
				Description: "Pass credentials in WinRM session to create a System.Management.Automation.PSCredential. (default: false, environment variable: AD_WINRM_PASS_CREDENTIALS)",
			},
			"domain_controller": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_DC", ""),
				Description: "Use a specific domain controller. (default: none, environment variable: AD_DC)",
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

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {
	cfg, err := config.NewConfig(d)
	if err != nil {
		return nil, err
	}
	pcfg := config.NewProviderConf(cfg)
	return pcfg, nil
}

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	// k is ignored here, but wee need to include it in the function's
	// signature in order to match the one defined for DiffSuppressFunc
	return strings.EqualFold(old, new)
}
