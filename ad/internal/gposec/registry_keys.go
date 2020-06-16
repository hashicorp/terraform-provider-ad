package gposec

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ini.v1"
)

// RegistryKeys represents the Registry Keys section of the Security Settings GPO extension
type RegistryKeys struct {
	Keys []string
}

//SetResourceData populates the resource's filed for the given section using the struct's data.
func (r *RegistryKeys) SetResourceData(section string, d *schema.ResourceData) error {
	out := []map[string]interface{}{}
	for _, valuesLine := range r.Keys {
		values := strings.SplitN(valuesLine, ",", 3)

    if len(values) != 3 {
			return fmt.Errorf("invalid registry keys line: %s", valuesLine)
		}

		value := map[string]interface{}{
			"key_name":         values[0],
			"propagation_mode": values[1],
			"acl":              values[2],
		}
		out = append(out, value)
	}
	return d.Set(section, out)
}

//SetIniData populates the INI file with data from this struct
func (r *RegistryKeys) SetIniData(f *ini.File) error {
	if len(r.Keys) == 0 {
		return nil
	}
	sectionName := "Registry Keys"
	sectionBody := strings.Join(r.Keys, "\r\n")
	f.NewRawSection(sectionName, fmt.Sprintf("%s\r\n", sectionBody))
	return nil
}

// NewRegistryKeysFromResource returns a new struct based on the resoruce's values
func NewRegistryKeysFromResource(data interface{}) (IniSetSection, error) {
	out := &RegistryKeys{Keys: []string{}}
	for _, item := range data.(*schema.Set).List() {
		rk := item.(map[string]interface{})
		value := fmt.Sprintf(`"%s",%s,"%s"`, rk["key_name"].(string), rk["propagation_mode"].(string), rk["acl"].(string))
		out.Keys = append(out.Keys, value)
	}
	return out, nil
}

// LoadRegistryKeysFromIni loads the data from the related INI section inside the given SecuritySettings
// struct
func LoadRegistryKeysFromIni(sectionName string, iniFile *ini.File, cfg *SecuritySettings) error {
	section, err := iniFile.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while parsing section %q: %s", sectionName, err)
	}
	cfg.RegistryKeys = &RegistryKeys{Keys: section.KeyStrings()}

	return nil
}
