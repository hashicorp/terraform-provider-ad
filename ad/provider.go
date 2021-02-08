package ad

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

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
			"winrm_use_ntlm": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AD_WINRM_USE_NTLM", false),
				Description: "Use NTLM authentication. (default: false, environment variable: AD_WINRM_USE_NTLM)",
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
	Configuration  *ProviderConfig
	winRMClients   []*winrm.Client
	winRMCPClients []*winrmcp.Winrmcp
	mx             *sync.Mutex
}

// AcquireWinRMClient get a thread safe WinRM client from the pool. Create a new one if the pool is empty
func (pcfg ProviderConf) AcquireWinRMClient() (winRMClient *winrm.Client, err error) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	if len(pcfg.winRMClients) == 0 {
		winRMClient, err = GetWinRMConnection(*pcfg.Configuration)
		if err != nil {
			return nil, err
		}
	} else {
		winRMClient = pcfg.winRMClients[0]
		pcfg.winRMClients = pcfg.winRMClients[1:]
	}
	return winRMClient, nil
}

// ReleaseWinRMClient returns a thread safe WinRM client after usage to the pool.
func (pcfg ProviderConf) ReleaseWinRMClient(winRMClient *winrm.Client) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	pcfg.winRMClients = append(pcfg.winRMClients, winRMClient)
}

// AcquireWinRMCPClient get a thread safe WinRM client from the pool. Create a new one if the pool is empty
func (pcfg ProviderConf) AcquireWinRMCPClient() (winRMCPClient *winrmcp.Winrmcp, err error) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	if len(pcfg.winRMClients) == 0 {
		winRMCPClient, err = GetWinRMCPConnection(*pcfg.Configuration)
		if err != nil {
			return nil, err
		}
	} else {
		winRMCPClient = pcfg.winRMCPClients[0]
		pcfg.winRMCPClients = pcfg.winRMCPClients[1:]
	}
	return winRMCPClient, nil
}

// ReleaseWinRMCPClient returns a thread safe WinRM client after usage to the pool.
func (pcfg ProviderConf) ReleaseWinRMCPClient(winRMCPClient *winrmcp.Winrmcp) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	pcfg.winRMCPClients = append(pcfg.winRMCPClients, winRMCPClient)
}

// isConnectionTypeLocal check if connection is local
func (pcfg ProviderConf) isConnectionTypeLocal() bool {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()

	log.Printf("[DEBUG] Checking if connection should be local")
	isLocal := false
	if runtime.GOOS == "windows" {
		if pcfg.Configuration.WinRMHost == "" && pcfg.Configuration.WinRMUsername == "" && pcfg.Configuration.WinRMPassword == "" {
			log.Printf("[DEBUG] Matching criteria for local execution")
			isLocal = true
		}
	}
	log.Printf("[DEBUG] Local connection ? %t", isLocal)
	return isLocal
}

func initProviderConfig(d *schema.ResourceData) (interface{}, error) {

	cfg := NewConfig(d)

	pcfg := ProviderConf{
		Configuration:  &cfg,
		winRMClients:   make([]*winrm.Client, 0),
		winRMCPClients: make([]*winrmcp.Winrmcp, 0),
		mx:             &sync.Mutex{},
	}

	return pcfg, nil
}

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	// k is ignored here, but wee need to include it in the function's
	// signature in order to match the one defined for DiffSuppressFunc
	return strings.EqualFold(old, new)
}
