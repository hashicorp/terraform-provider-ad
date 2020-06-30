package gposec

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ini.v1"
)

// RestrictedGroup represents a group that has its membership attributes managed by a GPO
type RestrictedGroup struct {
	GroupName    string
	GroupMembers string
	GroupParents string
}

// RestrictedGroups represents the Restricted Groups section of the Security Settings GPO extension
type RestrictedGroups struct {
	Groups []RestrictedGroup `mapstructure:"omitempty"`
}

// SetResourceData populates resource data based on the RestrictedGroups field values
func (r *RestrictedGroups) SetResourceData(section string, d *schema.ResourceData) error {
	out := []map[string]interface{}{}
	for _, group := range r.Groups {
		grp := map[string]interface{}{
			"group_name":     group.GroupName,
			"group_members":  group.GroupMembers,
			"group_memberof": group.GroupParents,
		}
		out = append(out, grp)
	}
	return d.Set(section, out)
}

//SetIniData populates the INI file with data from this struct
func (r *RestrictedGroups) SetIniData(f *ini.File) error {
	if len(r.Groups) == 0 {
		return nil
	}
	sectionName := "Group Membership"
	section, err := f.NewSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while creation INI Section %q", sectionName)
	}

	for _, group := range r.Groups {
		_, err := section.NewKey(fmt.Sprintf("%s__Members", group.GroupName), group.GroupMembers)
		if err != nil {
			return fmt.Errorf("error while creating new key for members of group %q: %s", group.GroupName, err)
		}
		_, err = section.NewKey(fmt.Sprintf("%s__Memberof", group.GroupName), group.GroupParents)
		if err != nil {
			return fmt.Errorf("error while creating new key for parents of group %q: %s", group.GroupName, err)
		}
	}
	return nil
}

// NewRestrictedGroupsFromResource returns a new struct based on the resoruce's values
func NewRestrictedGroupsFromResource(data interface{}) (IniSetSection, error) {
	out := &RestrictedGroups{Groups: []RestrictedGroup{}}
	for _, item := range data.(*schema.Set).List() {
		rgs := item.(map[string]interface{})
		rg := RestrictedGroup{
			GroupName:    rgs["group_name"].(string),
			GroupMembers: rgs["group_members"].(string),
			GroupParents: rgs["group_memberof"].(string),
		}
		out.Groups = append(out.Groups, rg)
	}
	return out, nil
}

// LoadRestrictedGroupsFromIni loads the data from the related INI section inside the given SecuritySettings
// struct
func LoadRestrictedGroupsFromIni(sectionName string, iniFile *ini.File, cfg *SecuritySettings) error {
	section, err := iniFile.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("error while parsing section %q: %s", sectionName, err)
	}
	out := &RestrictedGroups{Groups: []RestrictedGroup{}}
	// First pass, gather group names
	groups := make(map[string]int)
	for _, k := range section.KeyStrings() {
		keyParts := strings.Split(k, "__")
		if len(keyParts) != 2 {
			return fmt.Errorf("invalid key while processing restricted groups: %q", k)
		}
		groups[keyParts[0]] = 1
	}

	//Second pass, populate structure
	for group := range groups {
		g := RestrictedGroup{GroupName: group}
		if v, err := section.GetKey(fmt.Sprintf("%s__Members", group)); err == nil {
			g.GroupMembers = v.Value()
		}
		if v, err := section.GetKey(fmt.Sprintf("%s__Memberof", group)); err == nil {
			g.GroupParents = v.Value()
		}
		out.Groups = append(out.Groups, g)
	}
	cfg.RestrictedGroups = out
	return nil
}
