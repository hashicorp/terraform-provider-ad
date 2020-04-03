package ldaphelper

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
)

// Models

// EncodePassword takes a string, puts it between quotes, and converts it to UTF-16LE
// This is all required by the AD specification
//
func EncodePassword(password string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	quotedPassword := fmt.Sprintf("\"%s\"", password)
	utfEncodedPwd, err := utf16.NewEncoder().String(quotedPassword)
	if err != nil {
		return "", err
	}
	return utfEncodedPwd, nil
}

// GetResultFromLDAP sends a query to the LDAP server and returns a result.
func GetResultFromLDAP(conn *ldap.Conn, filter, base string, scope int, attrs []string) ([]*ldap.Entry, error) {
	sr := ldap.NewSearchRequest(
		base, scope, ldap.NeverDerefAliases,
		0, 0, false,
		filter, attrs, nil)

	result, err := conn.Search(sr)
	if err != nil {
		return nil, fmt.Errorf("ldap search failed. Filter was %q, base was %q, scope was %d. error: %s", filter, base, scope, err)
	}

	return result.Entries, nil
}

// GetDefaultNamingContext will return the domain name of the forest our DC belongs to
func GetDefaultNamingContext(conn *ldap.Conn) (string, error) {
	entries, err := GetResultFromLDAP(conn, "(defaultNamingContext=*)", "", ldap.ScopeBaseObject, nil)
	if err != nil {
		return "", err
	}

	if len(entries) > 1 {
		return "nil", fmt.Errorf("multiple entries found when querying for the default naming context. Aborting")
	}
	if len(entries) < 1 {
		return "nil", fmt.Errorf("No entries while querying for the default naming context")
	}

	entry := entries[0]
	dnc := entry.GetAttributeValue("defaultNamingContext")

	return dnc, nil

}
