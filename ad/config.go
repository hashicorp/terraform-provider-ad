package ad

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/masterzen/winrm"
)

// ProviderConfig holds all the information necessary to configure the provider
type ProviderConfig struct {
	LDAPHost      string
	LDAPPort      int
	LDAPUsername  string
	LDAPPassword  string
	LDAPProtocol  string
	LDAPInsecure  bool
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMProto    string
	WinRMInsecure bool
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) ProviderConfig {
	// ldap
	ldapHost := d.Get("dc_hostname").(string)
	ldapPort := d.Get("dc_port").(int)
	ldapUsername := d.Get("bind_username").(string)
	ldapPassword := d.Get("bind_password").(string)
	ldapProtocol := d.Get("proto").(string)
	ldapInsecure := d.Get("allow_insecure_certs").(bool)
	// winRM
	winRMUsername := d.Get("winrm_username").(string)
	if winRMUsername == "" {
		winRMUsername = ldapUsername
	}

	winRMPassword := d.Get("winrm_password").(string)
	if winRMPassword == "" {
		winRMPassword = ldapPassword
	}

	winRMHost := d.Get("winrm_hostname").(string)
	if winRMHost == "" {
		winRMHost = ldapHost
	}
	winRMPort := d.Get("winrm_port").(int)
	winRMProto := d.Get("winrm_proto").(string)
	winRMInsecure := d.Get("winrm_insecure").(bool)

	cfg := ProviderConfig{
		LDAPHost:      ldapHost,
		LDAPPort:      ldapPort,
		LDAPUsername:  ldapUsername,
		LDAPPassword:  ldapPassword,
		LDAPProtocol:  ldapProtocol,
		LDAPInsecure:  ldapInsecure,
		WinRMHost:     winRMHost,
		WinRMPort:     winRMPort,
		WinRMProto:    winRMProto,
		WinRMUsername: winRMUsername,
		WinRMPassword: winRMPassword,
		WinRMInsecure: winRMInsecure,
	}

	return cfg
}

func buildURL(config ProviderConfig, rootDSE bool) string {
	protocol := config.LDAPProtocol
	hostname := config.LDAPHost
	port := config.LDAPPort
	if rootDSE {
		return fmt.Sprintf("%s://%s:%d/rootDSE", protocol, hostname, port)
	}
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

// GetLDAPConnection returns an LDAP connection
func GetLDAPConnection(config ProviderConfig, rootDSE bool) (*ldap.Conn, error) {
	ldapURL := buildURL(config, rootDSE)
	var conn *ldap.Conn
	var err error
	if config.LDAPProtocol == "ldap" {
		conn, err = ldap.DialURL(ldapURL)
	} else if strings.ToLower(config.LDAPProtocol) == "ldaps" {
		conn, err = ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: config.LDAPInsecure}))
	} else {
		return nil, fmt.Errorf("invalid protocol %q specified", config.LDAPProtocol)
	}
	if err != nil {
		return nil, err
	}

	err = conn.Bind(config.LDAPUsername, config.LDAPPassword)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetWinRMConnection returns a WinRM connection
func GetWinRMConnection(config ProviderConfig) (*winrm.Client, error) {
	useHTTPS := false
	if strings.ToLower(config.WinRMProto) == "https" {
		useHTTPS = true
	}
	if config.WinRMHost == "" {
		config.WinRMHost = config.LDAPHost
	}

	if config.WinRMUsername == "" {
		config.WinRMUsername = config.LDAPUsername
	}

	if config.WinRMPassword == "" {
		config.WinRMPassword = config.LDAPPassword
	}

	endpoint := winrm.NewEndpoint(config.WinRMHost, config.WinRMPort, useHTTPS,
		config.WinRMInsecure, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, config.WinRMUsername, config.WinRMPassword)
	if err != nil {
		return nil, err
	}

	return client, nil
}
