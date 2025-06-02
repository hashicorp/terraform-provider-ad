// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
	"gopkg.in/ini.v1"
)

func TestRestrictedGroupsSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rk := &RestrictedGroups{
		Groups: []RestrictedGroup{
			{
				GroupName:    "group1",
				GroupMembers: "group2,group4",
				GroupParents: "group3",
			},
			{
				GroupName:    "group2",
				GroupMembers: "",
				GroupParents: "group1",
			},
			{
				GroupName:    "group3",
				GroupMembers: "group1",
				GroupParents: "",
			},
		},
	}

	err := rk.SetResourceData("restricted_groups", d)
	if err != nil {
		t.Error(err)
	}

	if d.Get("restricted_groups") == nil {
		t.Error("restricted_groups set is nil")
		t.FailNow()
	}

	rrSet := d.Get("restricted_groups").(*schema.Set)
	if rrSet.Len() == 0 {
		t.Error("empty restricted_groups set")
		t.FailNow()
	}

	for _, rg := range rrSet.List() {
		group := rg.(map[string]interface{})
		groupName := group["group_name"]
		groupMembers := group["group_members"]
		groupParents := group["group_memberof"]
		switch groupName {
		case "group1":
			if groupMembers != "group2,group4" || groupParents != "group3" {
				t.Errorf("group data are wrong for group %q", groupName)
			}
		case "group2":
			if groupMembers != "" || groupParents != "group1" {
				t.Errorf("group data are wrong for group %q", groupName)
			}
		case "group3":
			if groupMembers != "group1" || groupParents != "" {
				t.Errorf("group data are wrong for group %q", groupName)
			}
		default:
			t.Errorf("unexpected group name %q", groupName)
		}
	}
}

func newRGFromResource() (*RestrictedGroups, error) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rData := []map[string]interface{}{
		{
			"group_name":     "group1",
			"group_members":  "group2",
			"group_memberof": "group3",
		},
	}
	err := d.Set("restricted_groups", rData)
	if err != nil {
		return nil, err
	}

	rvSection, err := NewRestrictedGroupsFromResource(d.Get("restricted_groups"))
	if err != nil {
		return nil, err
	}
	rv := rvSection.(*RestrictedGroups)

	return rv, nil
}

func TestNewRestrictedGroupsFromResource(t *testing.T) {
	rv, err := newRGFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(rv.Groups) == 0 {
		t.Errorf("empty RestrictedGroups struct found")
		t.FailNow()
	}

	grp := rv.Groups[0]
	if grp.GroupName != "group1" || grp.GroupMembers != "group2" || grp.GroupParents != "group3" {
		t.Errorf(`RestrictedGroup structure did not contain expected values. `+
			`expected: name: group1, members: group2, parents: group3 `+
			`got: name: %s, members: %s, parents: %s`,
			grp.GroupName, grp.GroupMembers, grp.GroupParents)
	}
}

func TestRestrictedGroupsSetIniData(t *testing.T) {
	rv, err := newRGFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}
	iniFile := ini.Empty(loadOpts)
	ini.LineBreak = "\r\n"

	err = rv.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section := iniFile.Section("Group Membership")
	gm, err := section.GetKey(fmt.Sprintf("%s__Members", rv.Groups[0].GroupName))
	if err != nil {
		t.Errorf(fmt.Sprintf("key group1__Members wasn't found: %s", err))
	}
	if gm != nil && gm.Value() != "group2" {
		t.Errorf(fmt.Sprintf("unexpected value for group1__Members. expected group2 got %s", gm.Value()))
	}

	gm, err = section.GetKey(fmt.Sprintf("%s__Memberof", rv.Groups[0].GroupName))
	if err != nil {
		t.Errorf(fmt.Sprintf("key group1__Memberof wasn't found: %s", err))
	}
	if gm != nil && gm.Value() != "group3" {
		t.Errorf(fmt.Sprintf("unexpected value for group1__Memberof. expected group2 got %s", gm.Value()))
	}

	iniFile = ini.Empty(loadOpts)
	err = rv.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section = iniFile.Section("Restricted Groups")

	if section.Body() != "" {
		t.Errorf(`RestrictedGroups section body was not empty. actual value: %s`, section.Body())
	}
}

func TestLoadRestrictedGroupsFromIni(t *testing.T) {

	cfg := NewSecuritySettings()

	iniData := `
	[Restricted Groups]
	group1__Members=group2
	group1__Memberof=group3
	`

	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}
	iniFile, err := ini.LoadSources(loadOpts, []byte(iniData))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = LoadRestrictedGroupsFromIni("Restricted Groups", iniFile, cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(cfg.RestrictedGroups.Groups) != 1 {
		t.Errorf("Invalid number of Groups found in RestrictedGroups struct. Expected 1 got %d", len(cfg.RestrictedGroups.Groups))
		t.FailNow()
	}

	err = LoadRestrictedGroupsFromIni("Not Restricted Groups", iniFile, cfg)
	if !strings.Contains(err.Error(), "error while parsing section") {
		t.Error(err)
	}
}
