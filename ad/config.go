package ad

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jcmturner/gokrb5/v8/iana/etypeID"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/masterzen/winrm"
	"github.com/masterzen/winrm/soap"
	"github.com/packer-community/winrmcp/winrmcp"
)

// ProviderConfig holds all the information necessary to configure the provider
type ProviderConfig struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMProto    string
	WinRMInsecure bool
	KrbRealm      string
	KrbConfig     string
	KrbSpn        string
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) ProviderConfig {
	// winRM
	winRMUsername := d.Get("winrm_username").(string)
	winRMPassword := d.Get("winrm_password").(string)
	winRMHost := d.Get("winrm_hostname").(string)
	winRMPort := d.Get("winrm_port").(int)
	winRMProto := d.Get("winrm_proto").(string)
	winRMInsecure := d.Get("winrm_insecure").(bool)
	krbRealm := d.Get("krb_realm").(string)
	krbConfig := d.Get("krb_conf").(string)
	krbSpn := d.Get("krb_spn").(string)

	cfg := ProviderConfig{
		WinRMHost:     winRMHost,
		WinRMPort:     winRMPort,
		WinRMProto:    winRMProto,
		WinRMUsername: winRMUsername,
		WinRMPassword: winRMPassword,
		WinRMInsecure: winRMInsecure,
		KrbRealm:      krbRealm,
		KrbConfig:     krbConfig,
		KrbSpn:        krbSpn,
	}

	return cfg
}

// GetWinRMConnection returns a WinRM connection
func GetWinRMConnection(config ProviderConfig) (*winrm.Client, error) {
	useHTTPS := false
	if strings.ToLower(config.WinRMProto) == "https" {
		useHTTPS = true
	}

	endpoint := winrm.NewEndpoint(config.WinRMHost, config.WinRMPort, useHTTPS,
		config.WinRMInsecure, nil, nil, nil, 0)

	var winrmClient *winrm.Client
	var err error
	if config.KrbRealm != "" {
		params := winrm.DefaultParameters
		params.TransportDecorator = NewKerberosTransporter(config)
		winrmClient, err = winrm.NewClientWithParameters(endpoint, "", "", params)
	} else {
		winrmClient, err = winrm.NewClient(endpoint, config.WinRMUsername, config.WinRMPassword)
	}
	if err != nil {
		return nil, err
	}

	return winrmClient, nil
}

// GetWinRMCPConnection sets up a winrmcp client that can be used to upload files to the DC.
func GetWinRMCPConnection(config ProviderConfig) (*winrmcp.Winrmcp, error) {
	useHTTPS := false
	if config.WinRMProto == "https" {
		useHTTPS = true
	}
	addr := fmt.Sprintf("%s:%d", config.WinRMHost, config.WinRMPort)
	cfg := winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     config.WinRMUsername,
			Password: config.WinRMPassword,
		},
		Https:                 useHTTPS,
		Insecure:              config.WinRMInsecure,
		MaxOperationsPerShell: 15,
	}

	if config.KrbRealm != "" {
		cfg.TransportDecorator = NewKerberosTransporter(config)
	}

	return winrmcp.New(addr, &cfg)
}

type KerberosTransporter struct {
	Username  string
	Password  string
	Domain    string
	Hostname  string
	Port      int
	SPN       string
	KrbConf   string
	transport *http.Transport
}

func NewKerberosTransporter(config ProviderConfig) func() winrm.Transporter {
	return func() winrm.Transporter {
		return &KerberosTransporter{
			Username: config.WinRMUsername,
			Password: config.WinRMPassword,
			Domain:   config.KrbRealm,
			Hostname: config.WinRMHost,
			Port:     config.WinRMPort,
			KrbConf:  config.KrbConfig,
			SPN:      config.KrbSpn,
		}
	}
}

func (c *KerberosTransporter) Transport(endpoint *winrm.Endpoint) error {
	dial := (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial

	proxyfunc := http.ProxyFromEnvironment

	transport := &http.Transport{
		Proxy: proxyfunc,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: endpoint.Insecure,
			ServerName:         endpoint.TLSServerName,
		},
		Dial:                  dial,
		ResponseHeaderTimeout: endpoint.Timeout,
	}

	c.transport = transport

	return nil
}

func (c *KerberosTransporter) Post(_ *winrm.Client, request *soap.SoapMessage) (string, error) {
	var cfg *config.Config
	if c.KrbConf != "" {
		loadedCfg, err := config.Load(c.KrbConf)
		if err != nil {
			return "", err
		}
		cfg = loadedCfg
	} else {
		cfg = config.New()
		cfg.LibDefaults.DNSLookupKDC = false
		cfg.LibDefaults.DNSLookupRealm = false
		cfg.LibDefaults.PermittedEnctypes = []string{"aes128-cts-hmac-sha1-96", "aes256-cts-hmac-sha1-96",
			"aes128-cts-hmac-sha256-128", "aes256-cts-hmac-sha384-192"}
		cfg.LibDefaults.DefaultTGSEnctypes = []string{"aes128-cts-hmac-sha1-96", "aes256-cts-hmac-sha1-96",
			"aes128-cts-hmac-sha256-128", "aes256-cts-hmac-sha384-192"}
		cfg.LibDefaults.DefaultTktEnctypes = []string{"aes128-cts-hmac-sha1-96", "aes256-cts-hmac-sha1-96",
			"aes128-cts-hmac-sha256-128", "aes256-cts-hmac-sha384-192"}
		cfg.LibDefaults.DefaultRealm = c.Domain
		cfg.LibDefaults.UDPPreferenceLimit = 1
		cfg.LibDefaults.PreferredPreauthTypes = []int{17, 16, 15, 14}

		var encTypeIds []int32
		for _, encType := range cfg.LibDefaults.PermittedEnctypes {
			encTypeIds = append(encTypeIds, etypeID.EtypeSupported(encType))
		}
		cfg.LibDefaults.PermittedEnctypeIDs = encTypeIds

		var dflTGSEncTypeIds []int32
		for _, encType := range cfg.LibDefaults.DefaultTGSEnctypes {
			dflTGSEncTypeIds = append(dflTGSEncTypeIds, etypeID.EtypeSupported(encType))
		}
		cfg.LibDefaults.DefaultTGSEnctypeIDs = dflTGSEncTypeIds

		var dflTKTEncTypeIds []int32
		for _, encType := range cfg.LibDefaults.DefaultTktEnctypes {
			dflTKTEncTypeIds = append(dflTKTEncTypeIds, etypeID.EtypeSupported(encType))
		}
		cfg.LibDefaults.DefaultTktEnctypeIDs = dflTKTEncTypeIds

		cfg.Realms = []config.Realm{
			{
				AdminServer:   []string{fmt.Sprintf("%s:749", c.Hostname)},
				KDC:           []string{fmt.Sprintf("%s:88", c.Hostname)},
				KPasswdServer: []string{c.Hostname},
				Realm:         c.Domain,
			},
		}

		cfg.DomainRealm = config.DomainRealm{
			c.Domain: c.Domain,
		}
	}

	// setup the kerberos client
	kerberosClient := client.NewWithPassword(c.Username, c.Domain, c.Password, cfg, client.DisablePAFXFAST(true),
		client.AssumePreAuthentication(true))

	// setup the spnego client using the kerberos client we got above
	spnegoCl := spnego.NewClient(kerberosClient, nil, c.SPN)

	if c.transport != nil {
		spnegoCl.Transport = c.transport
	}

	//create an http request
	winrmURL := fmt.Sprintf("http://%s:%d/wsman", c.Hostname, c.Port)
	winRMRequest, _ := http.NewRequest("POST", winrmURL, strings.NewReader(request.String()))
	winRMRequest.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")

	// Use the spnego client to make the http request
	resp, err := spnegoCl.Do(winRMRequest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http error while making kerberos authenticated winRM request: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
