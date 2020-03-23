package ldaphelper

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// User represents an AD User
type User struct {
	SAMAccountName string
	DisplayName    string
	Password       string
	PrincipalName  string
	UserContainer  string
	DomainDN       string
	ChangeAtLogin  bool
}

// BuildDN returns a User's DN
func (u *User) BuildDN() string {
	dn := fmt.Sprintf("CN=%s,CN=%s,%s", u.DisplayName, u.UserContainer, u.DomainDN)
	return dn
}

// AddUser builds a user object from the resource data and sends it to the LDAP server
//in the form of an Add request.
func (u *User) AddUser(conn *ldap.Conn) (*string, error) {
	dn := u.BuildDN()
	log.Printf("Adding user with DN: %q", dn)

	unicodePwd, err := EncodePassword(u.Password)
	if err != nil {
		return nil, fmt.Errorf("password encoding failed: %s", err)
	}

	addReq := ldap.NewAddRequest(dn, []ldap.Control{})
	addReq.Attribute("distinguishedName", []string{dn})
	addReq.Attribute("sAMAccountName", []string{u.SAMAccountName})
	addReq.Attribute("displayName", []string{u.DisplayName})
	addReq.Attribute("userPrincipalName", []string{u.PrincipalName})
	addReq.Attribute("unicodePwd", []string{unicodePwd})
	addReq.Attribute("userAccountControl", []string{"512"})
	addReq.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	if u.ChangeAtLogin {
		addReq.Attribute("pwdLastSet", []string{"0"})
	}
	err = conn.Add(addReq)
	if err != nil {
		return nil, err
	}

	return &dn, nil
}

// ModifyUser updates the user LDAP object based on resource data
func (u *User) ModifyUser(d *schema.ResourceData, conn *ldap.Conn) error {
	dn := u.BuildDN()
	log.Printf("Modifying user with DN: %q", dn)
	keyMap := map[string]string{
		"sam_account_name": "sAMAccountName",
		"display_name":     "displayName",
		"initial_password": "unicodePwd",
		"principal_name":   "userPrincipalName",
	}
	modReq := ldap.NewModifyRequest(dn, []ldap.Control{})

	for _, k := range []string{"sam_account_name", "display_name", "principal_name"} {
		if d.HasChange(k) {
			_, newVal := d.GetChange(k)
			ldapKey := keyMap[k]
			value := newVal.(string)
			modReq.Replace(ldapKey, []string{value})
		}
	}

	if d.HasChange("initial_password") {
		unicodePwd, err := EncodePassword(u.Password)
		if err != nil {
			return fmt.Errorf("password encoding failed: %s", err)
		}
		modReq.Replace("unicodePwd", []string{string(unicodePwd)})
	}

	if d.HasChange("change_at_next_login") {
		_, newVal := d.Get("change_at_next_login").(bool)
		if newVal {
			modReq.Replace("pwdLastSet", []string{"0"})
		} else {
			modReq.Delete("pwdLastSet", []string{"0"})
		}
	}
	err := conn.Modify(modReq)
	return err
}

// GetUserFromResource returns a user struct built from Resource data
func GetUserFromResource(d *schema.ResourceData) *User {
	u := User{
		SAMAccountName: d.Get("sam_account_name").(string),
		DisplayName:    d.Get("display_name").(string),
		Password:       d.Get("initial_password").(string),
		PrincipalName:  d.Get("principal_name").(string),
		DomainDN:       d.Get("domain_dn").(string),
		ChangeAtLogin:  d.Get("change_at_next_login").(bool),
		UserContainer:  "Users",
	}

	return &u
}

// GetUserFromLDAP returns a User struct based on data
// retrieved from the LDAP server.
func GetUserFromLDAP(conn *ldap.Conn, dn, domainDN string) (*User, error) {
	filter := fmt.Sprintf("(&(distinguishedName=%s)(objectClass=user))", dn)
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
	pwdLastSet := entry.GetAttributeValue("pwdLastSet")
	var changePassword bool
	if pwdLastSet != "0" {
		changePassword = false
	} else {
		changePassword = true
	}

	u := &User{
		SAMAccountName: entry.GetAttributeValue("sAMAccountName"),
		DisplayName:    entry.GetAttributeValue("displayName"),
		PrincipalName:  entry.GetAttributeValue("userPrincipalName"),
		ChangeAtLogin:  changePassword,
		UserContainer:  "Users",
	}

	return u, nil
}
