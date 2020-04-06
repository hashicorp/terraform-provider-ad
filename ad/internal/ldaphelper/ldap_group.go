package ldaphelper

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var scopeMap = map[string]int64{
	"global":    0x00000002,
	"local":     0x00000004,
	"universal": 0x00000008,
}

var typeMap = map[string]int64{
	"system":    0x00000001,
	"app_basic": 0x00000010,
	"app_query": 0x00000020,
	"security":  0x80000000,
}

// TranslateLDAPGroupType will translate the groupType LDAP attribute's value to a pair of strings.
// The first string is the group's type and the second string is the group's scope.
// https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-server-2003/cc755692(v=ws.10)
// https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-server-2003/cc781446(v=ws.10)
func TranslateLDAPGroupType(gval int64) (string, string, error) {
	var groupType string
	var groupScope string

	for k, v := range typeMap {
		if (gval & v) != 0 {
			groupType = k
			break
		}
	}

	for k, v := range scopeMap {
		if (gval & v) != 0 {
			groupScope = k
			break
		}
	}

	if groupType == "" || groupScope == "" {
		return "", "", fmt.Errorf("could not translate %d to a meaningful type and scope", gval)
	}
	return groupScope, groupType, nil
}

// GetLDAPGroupType returns an int64 representing
// a group's scope and type based on resource configuration
// https://docs.microsoft.com/en-gb/windows/win32/adschema/a-grouptype
func GetLDAPGroupType(groupScope, groupType string) (int64, error) {

	groupInt, ok := scopeMap[groupScope]
	if !ok {
		return -1, fmt.Errorf("invalid group scope %q", groupInt)
	}

	typeInt, ok := typeMap[groupType]
	if !ok {
		return -1, fmt.Errorf("invalid group type %q", groupType)
	}

	result := 0x0 | typeInt | groupInt
	return result, nil
}

// Group represents an AD Group
type Group struct {
	SAMAccountName string
	Name           string
	Container      string
	DomainDN       string
	Scope          string
	Type           string
}

// BuildDN returns a Group's DN
func (g *Group) BuildDN() string {
	if g.Container != "" {
		return fmt.Sprintf("CN=%s%s,%s", g.Name, g.Container, g.DomainDN)
	}
	return fmt.Sprintf("CN=%s,%s", g.Name, g.DomainDN)
}

// AddGroup builds a Group object from the resource data and sends it to the LDAP server
//in the form of an Add request.
func (g *Group) AddGroup(conn *ldap.Conn) (*string, error) {
	dn := g.BuildDN()
	log.Printf("Adding Group with DN: %q", dn)
	groupType, err := GetLDAPGroupType(g.Scope, g.Type)
	if err != nil {
		return nil, err
	}

	gt := strconv.FormatInt(groupType, 10)
	addReq := ldap.NewAddRequest(dn, []ldap.Control{})
	addReq.Attribute("distinguishedName", []string{dn})
	addReq.Attribute("sAMAccountName", []string{g.SAMAccountName})
	addReq.Attribute("cn", []string{g.Name})
	addReq.Attribute("instanceType", []string{"4"})
	addReq.Attribute("objectClass", []string{"top", "group"})
	addReq.Attribute("groupType", []string{gt})

	err = conn.Add(addReq)
	if err != nil {
		return nil, err
	}

	return &dn, nil
}

// ModifyGroup updates the Group LDAP object based on resource data
func (g *Group) ModifyGroup(d *schema.ResourceData, conn *ldap.Conn) error {
	dn := g.BuildDN()
	log.Printf("Modifying Group with DN: %q", dn)
	keyMap := map[string]string{
		"sam_account_name": "sAMAccountName",
		"display_name":     "cn",
	}
	modReq := ldap.NewModifyRequest(dn, []ldap.Control{})

	for _, k := range []string{"sam_account_name", "display_name"} {
		if d.HasChange(k) {
			_, newVal := d.GetChange(k)
			ldapKey := keyMap[k]
			value := newVal.(string)
			modReq.Replace(ldapKey, []string{value})
		}
	}

	if d.HasChange("scope") || d.HasChange("type") {
		groupType, err := GetLDAPGroupType(g.Scope, g.Type)
		if err != nil {
			return err
		}
		modReq.Replace("groupType", []string{fmt.Sprintf("%d", groupType)})
	}

	err := conn.Modify(modReq)
	return err
}

// GetGroupFromResource returns a Group struct built from Resource data
func GetGroupFromResource(d *schema.ResourceData) *Group {
	g := Group{
		Name:           d.Get("display_name").(string),
		SAMAccountName: d.Get("sam_account_name").(string),
		DomainDN:       d.Get("domain_dn").(string),
		// Container:      "Groups",
		Scope: d.Get("scope").(string),
		Type:  d.Get("type").(string),
	}

	return &g
}

// GetGroupFromLDAP returns a Group struct based on data
// retrieved from the LDAP server.
func GetGroupFromLDAP(conn *ldap.Conn, dn, domainDN string) (*Group, error) {
	filter := fmt.Sprintf("(&(distinguishedName=%s)(objectClass=group))", dn)
	entries, err := GetResultFromLDAP(conn, filter, domainDN, ldap.ScopeWholeSubtree, nil)
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

	lgt := entry.GetAttributeValue("groupType")
	ldapGroupType, err := strconv.ParseInt(lgt, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ldap groupType (%s) to int64: %s", lgt, err)
	}
	groupScope, groupType, err := TranslateLDAPGroupType(int64(ldapGroupType))
	if err != nil {
		return nil, err
	}

	g := &Group{
		SAMAccountName: entry.GetAttributeValue("sAMAccountName"),
		Name:           entry.GetAttributeValue("cn"),
		// Container:      "Groups",
		DomainDN: domainDN,
		Scope:    groupScope,
		Type:     groupType,
	}

	return g, nil
}
