package gposec

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/ini.v1"
)

// Unicode is a required section
type Unicode struct {
	Unicode string
}

// Version is a required section
type Version struct {
	Signature string `ini:"signature"`
	Revision  int
}

// SystemAccess is a header in the INF file that holds information for the two sections described above
type SystemAccess struct {
	*PasswordPolicies `ini:"System Access,omitempty,squash"`
	*AccountLockout   `ini:"System Access,omitempty,squash"`
}

// SecuritySettings is a data structure representing the contents of the security settings INF file.
// It has tags used to map both the contents of the INF file as well as the resource data.
type SecuritySettings struct {
	Unicode
	Version
	*SystemAccess     `ini:"System Access,omitempty" mapstructure:"system_access,omitempty"`
	*KerberosPolicy   `ini:"Kerberos Policy,omitempty" mapstructure:"kerberos_settings,omitempty"`
	*EventAudit       `ini:"Event Audit,omitempty" mapstructure:"event_audit_policy,omitempty"`
	*SystemLog        `ini:"System Log,omitempty" mapstructure:"system_log,omitempty"`
	*AuditLog         `ini:"Security Log,omitempty" mapstructure:"audit_log,omitempty"`
	*ApplicationLog   `ini:"Application Log,omitempty" mapstructure:"application_log,omitempty"`
	*RestrictedGroups `ini:"Group Membership,omitempty" mapstructure:"restricted_groups,omitempty"`
	*RegistryKeys     `ini:"Registry Keys,omitempty" mapstructure:"registry_keys,omitempty"`
	*RegistryValues   `ini:"Registry Values,omitempty" mapstructure:"registry_values,omitempty"`
	*SystemServices   `ini:"Service General Setting,omitempty" mapstructure:"system_services,omitempty"`
	*FileSystem       `ini:"File Security,omitempty" mapstructure:"filesystem,omitempty"`
}

// PopulateSecuritySettings populates the SecuritySettings struct from resource data
func (s *SecuritySettings) PopulateSecuritySettings(d *schema.ResourceData, iniFile *ini.File) error {
	for section, fn := range ListSectionGeneratorMap {
		keyFunc := fn.(func(interface{}, *SecuritySettings) error)
		l := d.Get(section).([]interface{})
		if len(l) == 0 || l[0] == nil {
			continue
		}
		// All TypeLists in the resource schema have MaxItems set to 1
		data := l[0].(map[string]interface{})

		err := keyFunc(data, s)
		if err != nil {
			return fmt.Errorf("failed while processing setting type %q, error: %s", section, err)
		}
	}

	err := iniFile.ReflectFrom(s)
	if err != nil {
		return err
	}

	for section, fn := range SetSectionGeneratorMap {
		keyFunc := fn.(func(interface{}) (IniSetSection, error))

		setSection, err := keyFunc(d.Get(section))
		if err != nil {
			return fmt.Errorf("failed while processing setting type %q, error: %s", section, err)
		}
		err = setSection.SetIniData(iniFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewSecuritySettings returns a SecuritySettings struct with the header already populated.
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/fa15485d-ae9f-456e-a08f-81f2e5725a7e
func NewSecuritySettings() *SecuritySettings {
	sec := &SecuritySettings{
		Unicode: Unicode{Unicode: "yes"},
		Version: Version{Signature: "\"$CHICAGO$\"", Revision: 1},
	}
	return sec
}

// iniListSection is used when we need to treat all List-typed schema elements the same way
type iniListSection interface {
	SetResourceData(string, *schema.ResourceData) error
}

// IniSetSection is used when we need to treat all Set-typed schema elements the same way
type IniSetSection interface {
	SetIniData(*ini.File) error
}

// GetSectionData returns one of SecuritySettings' nested structures based on the key
// provided
func (s *SecuritySettings) GetSectionData(section string, d *schema.ResourceData) error {
	if section == "gpo_container" {
		// Nothing to do here
		return nil
	}

	var iniSection iniListSection

	switch section {
	case "password_policies":
		iniSection = s.SystemAccess.PasswordPolicies
	case "account_lockout":
		iniSection = s.SystemAccess.AccountLockout
	case "kerberos_policy":
		iniSection = s.KerberosPolicy
	case "system_log":
		iniSection = &s.SystemLog.EventLogPolicy
	case "audit_log":
		iniSection = &s.AuditLog.EventLogPolicy
	case "application_log":
		iniSection = &s.ApplicationLog.EventLogPolicy
	case "event_audit":
		iniSection = s.EventAudit
	case "restricted_groups":
		iniSection = s.RestrictedGroups
	case "registry_values":
		iniSection = s.RegistryValues
	case "system_services":
		iniSection = s.SystemServices
	case "registry_keys":
		iniSection = s.RegistryKeys
	case "filesystem":
		iniSection = s.FileSystem
	default:
		return fmt.Errorf("key %q is unknown", section)
	}

	// check if structure is empty
	emptySection := reflect.New(reflect.TypeOf(iniSection).Elem()).Interface()
	if reflect.DeepEqual(emptySection, iniSection) {
		log.Printf("[DEBUG] section %q is empty", section)
		return nil
	}

	err := iniSection.SetResourceData(section, d)
	return err
}

// ListSectionGeneratorMap maps a schema name to a function that populates the corresponding
// SecuritySettings fields with resource data.
var ListSectionGeneratorMap = map[string]interface{}{
	"password_policies": WritePasswordPolicies,
	"account_lockout":   WriteAccountLockout,
	"kerberos_policy":   WriteKerberosPolicy,
	"system_log":        WriteSystemLog,
	"audit_log":         WriteAuditLog,
	"application_log":   WriteApplicationLog,
	"event_audit":       WriteEventAudit,
}

// SetSectionGeneratorMap maps a schema name to a function that returns an INI section from resource data
// The difference with the map above is that this one deals with schema elements that are Sets instead
// of Lists and therefore require different handling.
var SetSectionGeneratorMap = map[string]interface{}{
	"restricted_groups": NewRestrictedGroupsFromResource,
	"registry_values":   NewRegistryValuesFromResource,
	"system_services":   NewSystemServicesFromResource,
	"registry_keys":     NewRegistryKeysFromResource,
	"filesystem":        NewFileSystemFromResource,
}

// SetSectionParserMap maps INI section names to functions that parse the sections and populate
// the relevant SecuritySettings fields. The sections not included in this map are handled
// by ini.MapTo().
var SetSectionParserMap = map[string]interface{}{
	"Service General Setting": LoadSystemServicesFromIni,
	"Group Membership":        LoadRestrictedGroupsFromIni,
	"Registry Keys":           LoadRegistryKeysFromIni,
	"Registry Values":         LoadRegistryValuesFromIni,
	"File Security":           LoadFileSystemFromIni,
}

// Most of the schema blocks in the resource's config are items of type List
// with max size of 1. These can be just converted via mapstructure.
func genericSetResourceData(section string, data interface{}, d *schema.ResourceData) error {
	out := make(map[string]interface{})
	err := mapstructure.Decode(data, &out)
	if err != nil {
		return fmt.Errorf("error in genericSetResourceData: %s", err)
	}

	return d.Set(section, []map[string]interface{}{out})
}

// UTFEncodeIniFile returs a byte array containing the encoded version of a string.
// The string is encoded to UTF16-LE with Byte Order Mark.
func UTFEncodeIniFile(iniFile *ini.File) (*[]byte, error) {
	encoding := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	outWriter := bytes.NewBuffer([]byte{})
	utf16Writer := encoding.NewEncoder().Writer(outWriter)
	ini.LineBreak = "\r\n"
	_, err := iniFile.WriteTo(utf16Writer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload to UTF16-LE with BOM, error: %s", err)
	}

	outBuf := outWriter.Bytes()
	return &outBuf, nil
}

// ParseIniFile decodes the INF file and returns an IniFile populated
// with the data found in it. If ut16fDecode is true then it translates
// contents from UTF16.
func ParseIniFile(iniBytes []byte, utf16Decode bool) (*SecuritySettings, error) {
	cfg := &SecuritySettings{}
	buf := bytes.NewBuffer(iniBytes)
	var reader io.Reader
	if utf16Decode {
		encoding := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		reader = encoding.NewDecoder().Reader(buf)
	} else {
		reader = buf
	}

	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}

	f, err := ini.LoadSources(loadOpts, reader)
	if err != nil {
		return nil, fmt.Errorf("error while loading ini contents: %s", err)
	}

	err = f.MapTo(cfg)
	if err != nil {
		return nil, err
	}

	for section, fn := range SetSectionParserMap {
		if _, err := f.GetSection(section); err != nil {
			continue
		}
		keyFunc := fn.(func(string, *ini.File, *SecuritySettings) error)
		err := keyFunc(section, f, cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, err
}

// HandleSectionRead handles all the logic behind the provider's Read() method.
// For the purposes of the function below:
// "section": is one of the blocks in the resource's configuration.
// "hostData": is the golang structure representing the data we parsed from the
//             .inf file we downladed from the host
func HandleSectionRead(schemaKeys []string, hostData *SecuritySettings, d *schema.ResourceData) error {
	// each of the schemaKeys elements, apart from gpo_container, represent a block
	// in the resource's configuration. Each of these blocks is mapped to a section
	// in the .INF file.
	for _, section := range schemaKeys {

		// We get the a string representation of the section's data
		// and update the state with these values.
		// sectionData, err := hostData.GetSectionData(section, d)
		err := hostData.GetSectionData(section, d)
		if err != nil {
			return err
		}
	}
	return nil
}
