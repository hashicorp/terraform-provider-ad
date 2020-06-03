package gposec

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ini.v1"
)

// RegistryValues is used to popullate the Registry Values section of the inf file that is used by
// many GPO features to set values in the registry
type RegistryValues struct {
	Values []string
}

//SetResourceData populates the resource's filed for the given section using the struct's data.
func (r *RegistryValues) SetResourceData(section string, d *schema.ResourceData) error {
	out := []map[string]interface{}{}
	for _, valuesLine := range r.Values {
		values := strings.SplitN(valuesLine, ",", 3)
		value := map[string]interface{}{
			"key_name":   values[0],
			"value_type": values[1],
			"value":      values[2],
		}
		out = append(out, value)
	}
	return d.Set(section, out)
}

//SetIniData populates the INI file with data from this struct
func (r *RegistryValues) SetIniData(f *ini.File) error {
	if len(r.Values) == 0 {
		return nil
	}
	sectionName := "Registry Values"
	sectionBody := strings.Join(r.Values, "\r\n")
	f.NewRawSection(sectionName, fmt.Sprintf("%s\r\n", sectionBody))
	return nil
}

// NewRegistryValuesFromResource returns a new struct based on the resoruce's values
func NewRegistryValuesFromResource(data interface{}) (IniSetSection, error) {
	out := &RegistryValues{Values: []string{}}
	for _, item := range data.(*schema.Set).List() {
		rv := item.(map[string]interface{})
		value := fmt.Sprintf(`"%s",%s,"%s"`, rv["key_name"].(string), rv["value_type"].(string), rv["value"].(string))
		out.Values = append(out.Values, value)
	}
	return out, nil
}

// LoadRegistryValuesFromIni loads the data from the related INI section inside the given SecuritySettings
// struct
func LoadRegistryValuesFromIni(sectionName string, iniFile *ini.File, cfg *SecuritySettings) error {
	section, err := iniFile.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while parsing section %q: %s", sectionName, err)
	}
	cfg.RegistryValues = &RegistryValues{Values: section.KeyStrings()}

	return nil
}
