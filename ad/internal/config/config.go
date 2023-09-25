// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jcmturner/gokrb5/v8/iana/etypeID"
	"github.com/jcmturner/gokrb5/v8/keytab"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/masterzen/winrm"
	"github.com/masterzen/winrm/soap"
	"github.com/packer-community/winrmcp/winrmcp"
)

// Settings holds all the information necessary to configure the provider
type Settings struct {
	WinRMUsername        string
	WinRMPassword        string
	WinRMHost            string
	WinRMPort            int
	WinRMProto           string
	WinRMInsecure        bool
	KrbRealm             string
	KrbConfig            string
	KrbKeytab            string
	KrbSpn               string
	WinRMUseNTLM         bool
	WinRMPassCredentials bool
	DomainName           string
	DomainController     string
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) (*Settings, error) {
	// winRM
	winRMUsername := d.Get("winrm_username").(string)
	winRMPassword := d.Get("winrm_password").(string)
	winRMHost := d.Get("winrm_hostname").(string)
	winRMPort := d.Get("winrm_port").(int)
	winRMProto := d.Get("winrm_proto").(string)
	winRMInsecure := d.Get("winrm_insecure").(bool)
	krbRealm := d.Get("krb_realm").(string)
	krbConfig := d.Get("krb_conf").(string)
	krbKeytab := d.Get("krb_keytab").(string)
	krbSpn := d.Get("krb_spn").(string)
	winRMUseNTLM := d.Get("winrm_use_ntlm").(bool)
	winRMPassCredentials := d.Get("winrm_pass_credentials").(bool)
	domainController := d.Get("domain_controller").(string)

	cfg := &Settings{
		DomainName:           krbRealm,
		WinRMHost:            winRMHost,
		WinRMPort:            winRMPort,
		WinRMProto:           winRMProto,
		WinRMUsername:        winRMUsername,
		WinRMPassword:        winRMPassword,
		WinRMInsecure:        winRMInsecure,
		KrbRealm:             krbRealm,
		KrbConfig:            krbConfig,
		KrbKeytab:            krbKeytab,
		KrbSpn:               krbSpn,
		WinRMUseNTLM:         winRMUseNTLM,
		WinRMPassCredentials: winRMPassCredentials,
		DomainController:     domainController,
	}

	return cfg, nil
}

// GetWinRMConnection returns a WinRM connection
func GetWinRMConnection(settings *Settings) (*winrm.Client, error) {
	useHTTPS := false
	if strings.ToLower(settings.WinRMProto) == "https" {
		useHTTPS = true
	}

	endpoint := winrm.NewEndpoint(settings.WinRMHost, settings.WinRMPort, useHTTPS,
		settings.WinRMInsecure, nil, nil, nil, 0)

	var winrmClient *winrm.Client
	var err error
	if settings.KrbRealm != "" {
		params := winrm.DefaultParameters
		params.TransportDecorator = NewKerberosTransporter(settings)
		winrmClient, err = winrm.NewClientWithParameters(endpoint, "", "", params)
	} else {
		params := winrm.DefaultParameters
		if settings.WinRMUseNTLM {
			params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
		}
		winrmClient, err = winrm.NewClientWithParameters(endpoint, settings.WinRMUsername, settings.WinRMPassword, params)
	}

	if err != nil {
		return nil, err
	}

	return winrmClient, nil
}

// GetWinRMCPConnection sets up a winrmcp client that can be used to upload files to the DC.
func GetWinRMCPConnection(settings *Settings) (*winrmcp.Winrmcp, error) {
	useHTTPS := false
	if settings.WinRMProto == "https" {
		useHTTPS = true
	}
	addr := fmt.Sprintf("%s:%d", settings.WinRMHost, settings.WinRMPort)
	cfg := winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     settings.WinRMUsername,
			Password: settings.WinRMPassword,
		},
		Https:                 useHTTPS,
		Insecure:              settings.WinRMInsecure,
		MaxOperationsPerShell: 15,
	}

	if settings.KrbRealm != "" {
		cfg.TransportDecorator = NewKerberosTransporter(settings)
	}

	return winrmcp.New(addr, &cfg)
}

type KerberosTransporter struct {
	Username  string
	Password  string
	Domain    string
	Hostname  string
	Port      int
	Proto     string
	SPN       string
	KrbConf   string
	KrbKeytab string
	transport *http.Transport
}

func NewKerberosTransporter(settings *Settings) func() winrm.Transporter {
	return func() winrm.Transporter {
		return &KerberosTransporter{
			Username:  settings.WinRMUsername,
			Password:  settings.WinRMPassword,
			Domain:    settings.KrbRealm,
			Hostname:  settings.WinRMHost,
			Port:      settings.WinRMPort,
			Proto:     settings.WinRMProto,
			KrbConf:   settings.KrbConfig,
			KrbKeytab: settings.KrbKeytab,
			SPN:       settings.KrbSpn,
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
	var kerberosClient *client.Client
	if c.KrbKeytab != "" {
		cfg.LibDefaults.DefaultKeytabName = c.KrbKeytab
		keytab, err := keytab.Load(cfg.LibDefaults.DefaultKeytabName)
		if err != nil {
			return "", err
		}
		kerberosClient = client.NewWithKeytab(c.Username, c.Domain, keytab, cfg, client.DisablePAFXFAST(true),
			client.AssumePreAuthentication(true))
	} else {
		kerberosClient = client.NewWithPassword(c.Username, c.Domain, c.Password, cfg, client.DisablePAFXFAST(true),
			client.AssumePreAuthentication(true))
	}

	// setup the spnego client using the kerberos client we got above
	spnegoCl := spnego.NewClient(kerberosClient, nil, c.SPN)

	if c.transport != nil {
		spnegoCl.Transport = c.transport
	}

	//create an http request
	winrmURL := fmt.Sprintf("%s://%s:%d/wsman", c.Proto, c.Hostname, c.Port)
	winRMRequest, _ := http.NewRequest("POST", winrmURL, strings.NewReader(request.String()))
	winRMRequest.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")

	// Use the spnego client to make the http request
	resp, err := spnegoCl.Do(winRMRequest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var bodyMsg string
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			bodyMsg = fmt.Sprintf("Also there was an error while retrieving the response's body: %s", err)
		} else {
			bodyMsg = fmt.Sprintf("response body:\n%s", string(respBody))
		}
		return "", fmt.Errorf("http error while making kerberos authenticated winRM request: %d - %s. %s ", resp.StatusCode, resp.Status, bodyMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

// ProviderConf holds structures that are useful to the provider at runtime.
type ProviderConf struct {
	Settings       *Settings
	winRMClients   []*winrm.Client
	winRMCPClients []*winrmcp.Winrmcp
	mx             *sync.Mutex
}

func NewProviderConf(settings *Settings) *ProviderConf {
	pcfg := &ProviderConf{
		Settings:       settings,
		winRMClients:   make([]*winrm.Client, 0),
		winRMCPClients: make([]*winrmcp.Winrmcp, 0),
		mx:             &sync.Mutex{},
	}
	return pcfg
}

// AcquireWinRMClient get a thread safe WinRM client from the pool. Create a new one if the pool is empty
func (pcfg *ProviderConf) AcquireWinRMClient() (winRMClient *winrm.Client, err error) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	if len(pcfg.winRMClients) == 0 {
		winRMClient, err = GetWinRMConnection(pcfg.Settings)
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
func (pcfg *ProviderConf) ReleaseWinRMClient(winRMClient *winrm.Client) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	pcfg.winRMClients = append(pcfg.winRMClients, winRMClient)
}

// AcquireWinRMCPClient get a thread safe WinRM client from the pool. Create a new one if the pool is empty
func (pcfg *ProviderConf) AcquireWinRMCPClient() (winRMCPClient *winrmcp.Winrmcp, err error) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	if len(pcfg.winRMCPClients) == 0 {
		winRMCPClient, err = GetWinRMCPConnection(pcfg.Settings)
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
func (pcfg *ProviderConf) ReleaseWinRMCPClient(winRMCPClient *winrmcp.Winrmcp) {
	pcfg.mx.Lock()
	defer pcfg.mx.Unlock()
	pcfg.winRMCPClients = append(pcfg.winRMCPClients, winRMCPClient)
}

// IsConnectionTypeLocal check if connection is local
func (pcfg *ProviderConf) IsConnectionTypeLocal() bool {
	log.Printf("[DEBUG] Checking if connection should be local")
	isLocal := false
	if runtime.GOOS == "windows" {
		if pcfg.Settings.WinRMHost == "" && pcfg.Settings.WinRMUsername == "" && pcfg.Settings.WinRMPassword == "" {
			log.Printf("[DEBUG] Matching criteria for local execution")
			isLocal = true
		}
	}
	log.Printf("[DEBUG] Local connection ? %t", isLocal)
	return isLocal
}

// IsPassCredentialsEnabled check if credentials should be passed
// requires that https be enabled
func (pcfg *ProviderConf) IsPassCredentialsEnabled() bool {
	log.Printf("[DEBUG] Checking to see if credentials should be passed")
	isPassCredentialsEnabled := false
	if pcfg.Settings.WinRMProto == "https" && pcfg.Settings.WinRMPassCredentials {
		log.Printf("[DEBUG] Matching criteria for passing credenitals")
		isPassCredentialsEnabled = true
	}
	log.Printf("[DEBUG] Pass Credentials ? %t", isPassCredentialsEnabled)
	return isPassCredentialsEnabled
}

// If a
func (pcfg *ProviderConf) IdentifyDomainController() string {
	log.Printf("[DEBUG] Checking to see if a domain controller was specified.")
	if pcfg.Settings.DomainController != "" {
		log.Printf("[DEBUG] Using specified domain controller for PowerShell commands.")
		return pcfg.Settings.DomainController
	}
	log.Printf("[DEBUG] Using the domain name instead of a specific domain controller for PowerShell commands.")
	return pcfg.Settings.DomainName
}
