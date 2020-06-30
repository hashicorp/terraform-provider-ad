package gposec

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ini.v1"
)

// SystemServices represents the System Services section of the Security Settings GPO extension
type SystemServices struct {
	Services []string
}

// SetResourceData populates resource data based on the SystemServices field values
func (r *SystemServices) SetResourceData(section string, d *schema.ResourceData) error {
	out := []map[string]interface{}{}
	for _, svcLine := range r.Services {
		fields := strings.SplitN(svcLine, ",", 3)
		if len(fields) != 3 {
			return fmt.Errorf("invalid services line: %s", svcLine)
		}

    svc := map[string]interface{}{
			"service_name": fields[0],
			"startup_mode": fields[1],
			"acl":          fields[2],
		}
		out = append(out, svc)
	}
	return d.Set(section, out)
}

// SetIniData populates the INI file with data.
func (r *SystemServices) SetIniData(f *ini.File) error {
	if len(r.Services) == 0 {
		return nil
	}
	sectionName := "Service General Setting"
	sectionBody := strings.Join(r.Services, "\r\n")
	f.NewRawSection(sectionName, fmt.Sprintf("%s\r\n", sectionBody))
	return nil
}

// NewSystemServicesFromResource returns a new SystemServices structure populated
// with data from the resources.
func NewSystemServicesFromResource(data interface{}) (IniSetSection, error) {
	out := &SystemServices{Services: []string{}}
	for _, item := range data.(*schema.Set).List() {
		ss := item.(map[string]interface{})
		service := fmt.Sprintf(`"%s",%s,"%s"`, ss["service_name"].(string), ss["startup_mode"].(string), ss["acl"].(string))
		out.Services = append(out.Services, service)
	}
	return out, nil
}

// LoadSystemServicesFromIni updates the given SecuritySettings struct with data
// parsed from the INI file
func LoadSystemServicesFromIni(sectionName string, iniFile *ini.File, cfg *SecuritySettings) error {
	section, err := iniFile.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while parsing section %q: %s", sectionName, err)
	}
	keys := []string{}
	for _, k := range section.KeyStrings() {
		k = strings.ReplaceAll(k, "\"", "")
		keys = append(keys, k)
	}
	cfg.SystemServices = &SystemServices{Services: keys}

	return nil
}
