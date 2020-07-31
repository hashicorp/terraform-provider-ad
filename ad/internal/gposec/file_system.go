package gposec

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/ini.v1"
)

// FileSystem represents the File System section of the Security Settings GPO extension
type FileSystem struct {
	Paths []string
}

//SetResourceData populates the resource's filed for the given section using the struct's data.
func (r *FileSystem) SetResourceData(section string, d *schema.ResourceData) error {
	out := []map[string]interface{}{}
	for _, valuesLine := range r.Paths {
		values := strings.SplitN(valuesLine, ",", 3)
		value := map[string]interface{}{
			"path":             values[0],
			"propagation_mode": values[1],
			"acl":              values[2],
		}
		out = append(out, value)
	}
	//lintignore:R001
	return d.Set(section, out)
}

//SetIniData populates the INI file with data from this struct
func (r *FileSystem) SetIniData(f *ini.File) error {
	if len(r.Paths) == 0 {
		return nil
	}
	sectionName := "File Security"
	sectionBody := strings.Join(r.Paths, "\r\n")
	_, err := f.NewRawSection(sectionName, fmt.Sprintf("%s\r\n", sectionBody))
	if err != nil {
		return fmt.Errorf("error while setting section %q: %s", sectionName, err)
	}
	return nil
}

// NewFileSystemFromResource returns a new struct based on the resource's values
func NewFileSystemFromResource(data interface{}) (IniSetSection, error) {
	out := &FileSystem{Paths: []string{}}
	for _, item := range data.(*schema.Set).List() {
		fs := item.(map[string]interface{})
		value := fmt.Sprintf(`"%s",%s,"%s"`, fs["path"].(string), fs["propagation_mode"].(string), fs["acl"].(string))
		out.Paths = append(out.Paths, value)
	}
	return out, nil
}

// LoadFileSystemFromIni loads the data from the related INI section inside the given SecuritySettings
// struct
func LoadFileSystemFromIni(sectionName string, iniFile *ini.File, cfg *SecuritySettings) error {
	section, err := iniFile.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while parsing section %q: %s", sectionName, err)
	}
	cfg.FileSystem = &FileSystem{Paths: section.KeyStrings()}

	return nil
}
