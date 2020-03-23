package msad

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProviderConfig holds all the information necessary to configure the provider
type ProviderConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Protocol string
	Insecure bool
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) ProviderConfig {
	host := d.Get("dc_hostname").(string)
	port := d.Get("dc_port").(int)
	username := d.Get("bind_username").(string)
	password := d.Get("bind_password").(string)
	protocol := d.Get("proto").(string)
	insecure := d.Get("allow_insecure_certs").(bool)

	cfg := ProviderConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Protocol: protocol,
		Insecure: insecure,
	}

	return cfg
}

func buildURL(config ProviderConfig, rootDSE bool) string {
	protocol := config.Protocol
	hostname := config.Host
	port := config.Port
	if rootDSE {
		return fmt.Sprintf("%s://%s:%d/rootDSE", protocol, hostname, port)
	}
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

// GetConnection returns an LDAP connection
func GetConnection(config ProviderConfig, rootDSE bool) (*ldap.Conn, error) {
	ldapURL := buildURL(config, rootDSE)
	var conn *ldap.Conn
	var err error
	if config.Protocol == "ldap" {
		conn, err = ldap.DialURL(ldapURL)
	} else if config.Protocol == "ldaps" {
		conn, err = ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: config.Insecure}))
	} else {
		return nil, fmt.Errorf("invalid protocol %q specified", config.Protocol)
	}
	if err != nil {
		return nil, err
	}

	err = conn.Bind(config.Username, config.Password)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
