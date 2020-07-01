package ldaphelper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var accountControlMap = map[string]int64{
	"disabled":               0x00000002,
	"password_never_expires": 0x00010000,
}

// User represents an AD User
type User struct {
	SAMAccountName       string
	DisplayName          string
	Password             string
	PrincipalName        string
	UserContainer        string
	DomainDN             string
	Disabled             bool
	PasswordNeverExpires bool
}

// BuildDN returns a User's DN
func (u *User) BuildDN() string {
	dn := fmt.Sprintf("CN=%s,CN=%s,%s", u.DisplayName, u.UserContainer, u.DomainDN)
	return dn
}

// GetUserAccountControl returns the int64 value of userAccountControl that is used to set various
// flags for an AD user
// https://docs.microsoft.com/en-us/windows/win32/adschema/a-useraccountcontrol#remarks
func (u *User) GetUserAccountControl() int64 {
	var accountControlMap = map[string]int64{
		"disabled":               0x00000002,
		"password_never_expires": 0x00010000,
	}

	// Default value for users
	var userAccountControl int64 = 0x200

	if u.Disabled {
		userAccountControl = userAccountControl | accountControlMap["disabled"]
	}

	if u.PasswordNeverExpires {
		userAccountControl = userAccountControl | accountControlMap["password_never_expires"]
	}

	return userAccountControl
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

	uac := strconv.FormatInt(u.GetUserAccountControl(), 10)

	addReq := ldap.NewAddRequest(dn, []ldap.Control{})
	addReq.Attribute("distinguishedName", []string{dn})
	addReq.Attribute("sAMAccountName", []string{u.SAMAccountName})
	addReq.Attribute("displayName", []string{u.DisplayName})
	addReq.Attribute("userPrincipalName", []string{u.PrincipalName})
	addReq.Attribute("unicodePwd", []string{unicodePwd})
	addReq.Attribute("userAccountControl", []string{"512"})
	addReq.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	addReq.Attribute("userAccountControl", []string{uac})
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
			ldapKey := keyMap[k]
			value := d.Get(k).(string)
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

	if d.HasChange("disabled") || d.HasChange("password_never_expires") {
		uac := strconv.FormatInt(u.GetUserAccountControl(), 10)
		modReq.Replace("userAccountControl", []string{uac})
	}

	err := conn.Modify(modReq)
	return err
}

// GetUserFromResource returns a user struct built from Resource data
func GetUserFromResource(d *schema.ResourceData) *User {
	u := User{
		SAMAccountName:       d.Get("sam_account_name").(string),
		DisplayName:          d.Get("display_name").(string),
		Password:             d.Get("initial_password").(string),
		PrincipalName:        d.Get("principal_name").(string),
		DomainDN:             d.Get("domain_dn").(string),
		UserContainer:        "Users",
		Disabled:             d.Get("disabled").(bool),
		PasswordNeverExpires: d.Get("password_never_expires").(bool),
	}

	return &u
}

// GetUserFromLDAP returns a User struct based on data
// retrieved from the LDAP server.
func GetUserFromLDAP(conn *ldap.Conn, dn string) (*User, error) {
	filter := fmt.Sprintf("(&(distinguishedName=%s)(objectClass=user))", dn)
	domainDNIdx := strings.Index(dn, "dc=")
	domainDN := dn[domainDNIdx:]
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

	uac, err := strconv.ParseInt(entry.GetAttributeValue("userAccountControl"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error while parsing uac value from ldap (%q) into integer: %s", entry.GetAttributeValue("userAccountControl"), err)
	}
	disabled := uac&accountControlMap["disabled"] != 0
	passwordNeverExpires := uac&accountControlMap["password_never_expires"] != 0

	u := &User{
		SAMAccountName:       entry.GetAttributeValue("sAMAccountName"),
		DisplayName:          entry.GetAttributeValue("displayName"),
		PrincipalName:        entry.GetAttributeValue("userPrincipalName"),
		UserContainer:        "Users",
		Disabled:             disabled,
		PasswordNeverExpires: passwordNeverExpires,
		DomainDN:             domainDN,
	}

	return u, nil
}
