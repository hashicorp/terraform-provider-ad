package ldap

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-msad/msad"
)

func buildUrl(config *msad.Config) string {
	protocol := config.Protocol
	hostname := config.Host
	port := config.Port
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

// GetConnection returns an LDAP connection
func GetConnection(config *msad.Config) (*3.Conn, error) {
	ldapUrl := buildUrl(config)
	conn, err := v3.DialURL(ldapUrl)
	if err != nil {
		return nil, err 
	}
	if config.Protocol == "ldaps" {
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: config.Insecure})
	}
	return &conn, nil
}
