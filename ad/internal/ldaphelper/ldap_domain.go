package ldaphelper

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Domain represents an AD Domain
type Domain struct {
	DN          string
	NetbiosName string
	DomainName  string
}

// GetDomainFromResource returns a Domain struct built from Resource data
func GetDomainFromResource(d *schema.ResourceData) *Domain {
	domain := Domain{
		DN:          d.Get("dn").(string),
		NetbiosName: d.Get("netbios_name").(string),
		DomainName:  d.Get("domain_name").(string),
	}
	return &domain
}

// GetDomainFromLDAP returns a Domain struct based on data
// retrieved from the LDAP server
func GetDomainFromLDAP(dseConn *ldap.Conn, dn, netbiosName, domainName string) (*Domain, error) {

	dnc, err := GetDefaultNamingContext(dseConn)
	if err != nil {
		return nil, err
	}
	base := fmt.Sprintf("cn=partitions,cn=configuration,%s", dnc)

	filters := []string{}
	if dn != "" {
		filters = append(filters, fmt.Sprintf("(nCName=%s)", dn))
	}
	if netbiosName != "" {
		filters = append(filters, fmt.Sprintf("(nETBIOSName=%s)", netbiosName))
	}
	if domainName != "" {
		filters = append(filters, fmt.Sprintf("dnsRoot=%s)", domainName))
	}

	var filter string
	if len(filters) > 1 {
		filter = fmt.Sprintf("(&%s)", strings.Join(filters, ""))
	} else {
		filter = filters[0]
	}

	entries, err := GetResultFromLDAP(dseConn, filter, base, ldap.ScopeWholeSubtree, nil)
	if err != nil {
		return nil, err
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("multiple entries found for filter (%q). Aborting", filter)
	}
	if len(entries) < 1 {
		return nil, fmt.Errorf("No entries found for filter %q", filter)
	}

	entry := entries[0]
	domain := &Domain{
		DN:          entry.GetAttributeValue("nCName"),
		NetbiosName: entry.GetAttributeValue("nETBIOSName"),
		DomainName:  entry.GetAttributeValue("dnsRoot"),
	}

	return domain, nil
}
